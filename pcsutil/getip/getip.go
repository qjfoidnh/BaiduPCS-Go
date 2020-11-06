// Package getip 获取 ip 信息包
package getip

import (
	"github.com/iikira/BaiduPCS-Go/requester"
	"net"
	"net/http"
	"unsafe"
)

// IPInfoByClient 给定client获取ip地址
func IPInfoByClient(c *requester.HTTPClient) (ipAddr string, err error) {
	if c == nil {
		c = requester.NewHTTPClient()
	}

	body, err := c.Fetch(http.MethodGet, "https://api.ipify.org", nil, nil)
	if err != nil {
		return
	}

	ipAddr = *(*string)(unsafe.Pointer(&body))
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return "", ErrParseIP
	}
	return
}

//IPInfo 从ipify获取IP地址
func IPInfo(https bool) (ipAddr string, err error) {
	c := requester.NewHTTPClient()
	c.SetHTTPSecure(https)
	return IPInfoByClient(c)
}
