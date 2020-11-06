package panhome

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/expires"
	"github.com/iikira/BaiduPCS-Go/requester"
	"net/url"
)

const (
	// OperationSignature signature
	OperationSignature = "signature"
)

var (
	panBaiduComURL = &url.URL{
		Scheme: "https",
		Host:   "pan.baidu.com",
	}
	// PanHomeUserAgent PanHome User-Agent
	PanHomeUserAgent = "Mozilla/5.0"
)

type (
	PanHome struct {
		client *requester.HTTPClient
		ua     string
		bduss  string

		sign1, sign3 []rune
		timestamp    string

		signRes     SignRes
		signExpires expires.Expires
	}
)

func NewPanHome(client *requester.HTTPClient) *PanHome {
	ph := PanHome{}
	if client != nil {
		newC := *client
		ph.client = &newC
	}
	return &ph
}

func (ph *PanHome) lazyInit() {
	if ph.client == nil {
		ph.client = requester.NewHTTPClient()
	}
}
