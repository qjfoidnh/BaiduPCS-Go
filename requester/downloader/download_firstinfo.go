package downloader

import (
	"fmt"
	"net/http"
	"reflect"
)

type (
	DownloadFirstInfo struct {
		ContentLength int64
		ContentMD5    string
		ContentCRC32  string
		AcceptRanges  string
		Referer       string
	}
)

func NewDownloadFirstInfoByResp(contentLength int64, resp *http.Response) (dfi *DownloadFirstInfo) {
	dfi = &DownloadFirstInfo{}
	if resp == nil {
		dfi.ContentLength = contentLength
		return
	}
	if contentLength != resp.ContentLength {
		dfi.ContentLength = contentLength
	}
	dfi.AcceptRanges = resp.Header.Get("Accept-Ranges")
	dfi.Referer = resp.Header.Get("Referer")
	return
}

func (dfi *DownloadFirstInfo) Compare(n *DownloadFirstInfo) bool {
	if n == nil {
		return false
	}
	if dfi.ContentLength != n.ContentLength {
		return false
	}
	if dfi.AcceptRanges != n.AcceptRanges {
		return false
	}
	if dfi.Referer != n.Referer {
		return false
	}
	return true
}

// ToMap 转换为map
func (dfi *DownloadFirstInfo) ToMap() map[string]string {
	m := map[string]string{
		"Content-MD5":     dfi.ContentMD5,
		"x-bs-meta-crc32": dfi.ContentCRC32,
		"Accept-Ranges":   dfi.AcceptRanges,
		"Referer":         dfi.Referer,
	}
	return m
}

// ToMapByReflect 用reflect转换为map
func (dfi *DownloadFirstInfo) ToMapByReflect() map[string]string {
	te := reflect.TypeOf(dfi).Elem()
	ve := reflect.ValueOf(dfi).Elem()
	n := te.NumField()
	m := map[string]string{}
	for i := 0; i < n; i++ {
		f := te.Field(i)
		m[f.Name] = fmt.Sprint(ve.Field(i).Interface())
	}
	return m
}
