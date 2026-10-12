package requester

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires/cachemap"
	"github.com/rs/dnscache"
	mathrand "math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// MaxDuration 最大的Duration
	MaxDuration = 1<<63 - 1
)

var (
	localTCPAddrList = []*net.TCPAddr{}

	// preferIPv6 为真时，在域名同时具备 AAAA/A 记录时优先选择 IPv6 地址。
	preferIPv6 bool

	// ErrNoIPAddr 解析结果为空
	ErrNoIPAddr = errors.New("no IP addresses found")

	// ProxyAddr 代理地址
	ProxyAddr string

	// 指定走ProxyAddr的域名
	ProxyHostnameRules string

	// ErrProxyAddrEmpty 代理地址为空
	ErrProxyAddrEmpty = errors.New("proxy addr is empty")

	tcpCache = cachemap.GlobalCacheOpMap.LazyInitCachePoolOp("requester/tcp")

	// dnsResolver 通过 singleflight 合并同 host 查询；dnsLookupMu 串行化对标准库解析器的调用，
	// 避免并发 Lookup 时在 tryOneName 路径触发 panic。
	dnsResolver = &dnscache.Resolver{
		Timeout: 30 * time.Second,
	}
	dnsLookupMu sync.Mutex
)

// SetLocalTCPAddrList 设置网卡地址
func SetLocalTCPAddrList(ips ...string) {
	list := make([]*net.TCPAddr, 0, len(ips))
	for k := range ips {
		p := net.ParseIP(ips[k])
		if p == nil {
			continue
		}

		list = append(list, &net.TCPAddr{
			IP: p,
		})
	}
	localTCPAddrList = list
}

// SetPreferIPv6 设置是否优先使用 IPv6 地址。
//
// 本程序的 dnsResolver 依赖系统解析顺序(LookupHost 返回值)，而多数 Linux 的
// RFC3484 / gai.conf 规则会把 IPv4 排在前面，所以即使本机已有可用的 IPv6 出口，
// 实际仍会拨号到 IPv4。开启本项后，只要域名同时存在 AAAA 与 A 记录，就会在候选中
// 优先选择 IPv6 地址；若域名无 AAAA 记录，则自动回退到原有行为。
func SetPreferIPv6(b bool) {
	preferIPv6 = b
}

// pickAddr 按偏好从解析结果中选择地址；preferIPv6 为真时优先 IPv6。
func pickAddr(addrs []string) (string, error) {
	if len(addrs) == 0 {
		return "", ErrNoIPAddr
	}
	if !preferIPv6 {
		return addrs[0], nil
	}
	for _, a := range addrs {
		if ip := net.ParseIP(a); ip != nil && ip.To4() == nil {
			return a, nil
		}
	}
	return addrs[0], nil
}

func proxyFunc(req *http.Request) (*url.URL, error) {
	u, err := checkProxyAddr(ProxyAddr)
	if err != nil {
		return http.ProxyFromEnvironment(req)
	}

	if ProxyHostnameRules != "" {
		if strings.Contains(ProxyHostnameRules, req.URL.Hostname()) {
			return u, err
		} else {
			return http.ProxyFromEnvironment(req)
		}
	}
	return u, err
}

func getLocalTCPAddr() *net.TCPAddr {
	if len(localTCPAddrList) == 0 {
		return nil
	}
	i := mathrand.Intn(len(localTCPAddrList))
	return localTCPAddrList[i]
}

func getDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		LocalAddr: getLocalTCPAddr(),
		DualStack: true,
	}
}

func checkProxyAddr(proxyAddr string) (u *url.URL, err error) {
	if proxyAddr == "" {
		return nil, ErrProxyAddrEmpty
	}

	host, port, err := net.SplitHostPort(proxyAddr)
	if err == nil {
		u = &url.URL{
			Host: net.JoinHostPort(host, port),
		}
		return
	}

	u, err = url.Parse(proxyAddr)
	if err == nil {
		return
	}

	return
}

// SetGlobalProxy 设置代理
func SetGlobalProxy(proxyAddr string) {
	ProxyAddr = proxyAddr
}

// SetProxyHostnameRules 设置走代理的域名列表
func SetProxyHostnameRules(hostnames string) {
	ProxyHostnameRules = hostnames
}

// SetTCPHostBind 设置host绑定ip
func SetTCPHostBind(host, ip string) {
	tcpCache.Store(host, expires.NewDataExpires(net.ParseIP(ip), MaxDuration))
	return
}

func getServerName(address string) string {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}
	return host
}

// resolveTCPHost
// 解析的tcpaddr没有port!!!
func resolveTCPHost(ctx context.Context, host string) (ip net.IP, err error) {
	// host 已是字面量 IP 时直接返回，无需 DNS 查询
	if parsed := net.ParseIP(host); parsed != nil {
		return parsed, nil
	}

	// 标准库 DNS 偶发 panic 时转为 error，避免进程崩溃
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("dns lookup panic for %q: %v", host, r)
		}
	}()

	dnsLookupMu.Lock()
	addrs, lookupErr := dnsResolver.LookupHost(ctx, host)
	dnsLookupMu.Unlock()
	if lookupErr != nil {
		return nil, lookupErr
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no IP addresses found for host %q", host)
	}

	chosen, err := pickAddr(addrs)
	if err != nil {
		return nil, err
	}
	ip = net.ParseIP(chosen)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address %q for host %q", chosen, host)
	}
	return ip, nil
}

func dialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		host, portStr, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		data, err := cachemap.GlobalCacheOpMap.CacheOperationWithError("requester/tcp", host, func() (expires.DataExpires, error) {
			ip, err := resolveTCPHost(ctx, host)
			if err != nil {
				return nil, err
			}
			return expires.NewDataExpires(ip, 10*time.Minute), nil // 传值
		})
		if err != nil {
			return nil, err
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}

		return net.DialTCP(network, getLocalTCPAddr(), &net.TCPAddr{
			IP:   data.Data().(net.IP),
			Port: port, // 设置端口
		})
	}

	// 非 tcp 请求
	conn, err = getDialer().DialContext(ctx, network, address)
	return
}

func (h *HTTPClient) dialTLSFunc() func(network, address string) (tlsConn net.Conn, err error) {
	return func(network, address string) (tlsConn net.Conn, err error) {
		conn, err := dialContext(context.Background(), network, address)
		if err != nil {
			return nil, err
		}

		return tls.Client(conn, &tls.Config{
			ServerName:         getServerName(address),
			InsecureSkipVerify: !h.https,
		}), nil
	}
}
