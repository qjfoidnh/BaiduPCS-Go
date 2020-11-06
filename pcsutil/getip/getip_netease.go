package getip

import (
	"github.com/iikira/BaiduPCS-Go/pcsutil/jsonhelper"
	"github.com/iikira/BaiduPCS-Go/requester"
	"net"
	"net/http"
)

type (
	// IPResNetease 网易服务器获取ip返回的结果
	IPResNetease struct {
		Result  string `json:"result"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func IPInfoFromNeteaseByClient(c *requester.HTTPClient) (ipAddr string, err error) {
	resp, err := c.Req(http.MethodGet, "http://mam.netease.com/api/config/getClientIp", nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	res := &IPResNetease{}
	err = jsonhelper.UnmarshalData(resp.Body, res)
	if err != nil {
		return
	}

	ip := net.ParseIP(res.Result)
	if ip == nil {
		err = ErrParseIP
		return
	}

	ipAddr = res.Result
	return
}

// IPInfoFromNetease 从网易服务器获取ip
func IPInfoFromNetease() (ipAddr string, err error) {
	c := requester.NewHTTPClient()
	return IPInfoFromNeteaseByClient(c)
}
