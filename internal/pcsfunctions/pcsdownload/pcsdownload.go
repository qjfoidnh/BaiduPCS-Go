package pcsdownload

import (
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"net/http"
	"strconv"
)

// IsSkipMd5Checksum 是否忽略某些校验
func IsSkipMd5Checksum(size int64, md5Str string) bool {
	switch {
	case size == 1749504 && md5Str == "48bb9b0361dc9c672f3dc7b3ffcfde97": //8秒温馨提示
		fallthrough
	case size == 120 && md5Str == "6c1b84914588d09a6e5ec43605557457": //温馨提示文字版
		return true
	}
	return false
}

// BaiduPCSURLCheckFunc downloader 首次检查下载地址要执行的函数
func BaiduPCSURLCheckFunc(client *requester.HTTPClient, durl string) (contentLength int64, resp *http.Response, err error) {
	resp, err = client.Req(http.MethodGet, durl, nil, map[string]string{
		"Range": "bytes=0-" + strconv.FormatInt(baidupcs.InitRangeSize-1, 10),
	})
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return 0, nil, err
	}

	contentLengthStr := resp.Header.Get("x-bs-file-size")
	contentLength, _ = strconv.ParseInt(contentLengthStr, 10, 64)
	return contentLength, resp, nil
}
