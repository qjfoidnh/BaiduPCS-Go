package requester

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// HTTPClient http client
type HTTPClient struct {
	http.Client
	transport *http.Transport
	https     bool
	UserAgent string
}

type traceConn struct {
	net.Conn
	id string
}

func (c traceConn) Close() error {
	fmt.Printf("CONN[%s] Close() called, going to FIN -> TIME_WAIT", c.id)
	return c.Conn.Close()
}

// 拨号时包装
func dialTrace(ctx context.Context, network, addr string) (net.Conn, error) {
	d := net.Dialer{}
	raw, err := d.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s->%s@%d", raw.LocalAddr(), raw.RemoteAddr(), time.Now().UnixNano())
	fmt.Printf("CONN[%s] DialContext success", id)
	return traceConn{Conn: raw, id: id}, nil
}

// NewHTTPClient 返回 HTTPClient 的指针,
// 预设了一些配置
func NewHTTPClient() *HTTPClient {
	fmt.Println("new one")
	h := &HTTPClient{
		Client: http.Client{
			Timeout: 50 * time.Second,
			Transport: &http.Transport{
				DialContext:           dialTrace,
				MaxIdleConns:          20,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				// 关键：打开状态回调
				DisableKeepAlives: false,
			},
		},
		UserAgent: UserAgent,
	}
	h.Client.Jar, _ = cookiejar.New(nil)
	return h
}

func (h *HTTPClient) lazyInit() {
	if h.transport == nil {
		h.transport = &http.Transport{
			Proxy:       proxyFunc,
			DialContext: dialContext,
			Dial:        dial,
			// DialTLS:     h.dialTLSFunc(),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			TLSHandshakeTimeout:   20 * time.Second,
			DisableKeepAlives:     false,
			DisableCompression:    false, // gzip
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 25 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		}
		h.Client.Transport = h.transport
	}
}

// SetUserAgent 设置 UserAgent 浏览器标识
func (h *HTTPClient) SetUserAgent(ua string) {
	h.UserAgent = ua
}

// SetProxy 设置代理
func (h *HTTPClient) SetProxy(proxyAddr string) {
	h.lazyInit()
	u, err := checkProxyAddr(proxyAddr)
	if err != nil {
		h.transport.Proxy = http.ProxyFromEnvironment
		return
	}

	h.transport.Proxy = http.ProxyURL(u)
}

// SetCookiejar 设置 cookie
func (h *HTTPClient) SetCookiejar(jar http.CookieJar) {
	h.Client.Jar = jar
}

// ResetCookiejar 清空 cookie
func (h *HTTPClient) ResetCookiejar() {
	h.Jar, _ = cookiejar.New(nil)
}

// SetHTTPSecure 是否启用 https 安全检查, 强制不检查
func (h *HTTPClient) SetHTTPSecure(b bool) {
	h.https = b
	h.lazyInit()
	// 	if b {
	// 		h.transport.TLSClientConfig = nil
	// 	} else {
	h.transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	// 	}
}

// SetKeepAlive 设置 Keep-Alive
func (h *HTTPClient) SetKeepAlive(b bool) {
	h.lazyInit()
	h.transport.DisableKeepAlives = !b
}

// SetGzip 是否启用Gzip
func (h *HTTPClient) SetGzip(b bool) {
	h.lazyInit()
	h.transport.DisableCompression = !b
}

// SetResponseHeaderTimeout 设置目标服务器响应超时时间
func (h *HTTPClient) SetResponseHeaderTimeout(t time.Duration) {
	h.lazyInit()
	h.transport.ResponseHeaderTimeout = t
}

// SetTLSHandshakeTimeout 设置tls握手超时时间
func (h *HTTPClient) SetTLSHandshakeTimeout(t time.Duration) {
	h.lazyInit()
	h.transport.TLSHandshakeTimeout = t
}

// SetTimeout 设置 http 请求超时时间, 默认30s
func (h *HTTPClient) SetTimeout(t time.Duration) {
	h.Client.Timeout = t
}
