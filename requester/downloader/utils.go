package downloader

import (
	"github.com/iikira/BaiduPCS-Go/pcsverbose"
	"github.com/iikira/BaiduPCS-Go/requester"
	mathrand "math/rand"
	"mime"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"time"
)

var (
	// ContentRangeRE Content-Range 正则
	ContentRangeRE = regexp.MustCompile(`^.*? \d*?-\d*?/(\d*?)$`)

	// ranSource 随机数种子
	ranSource = mathrand.NewSource(time.Now().UnixNano())

	// ran 一个随机数实例
	ran = mathrand.New(ranSource)
)

// RandomNumber 生成指定区间随机数
func RandomNumber(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return ran.Intn(max-min) + min
}

// GetFileName 获取文件名
func GetFileName(uri string, client *requester.HTTPClient) (filename string, err error) {
	if client == nil {
		client = requester.NewHTTPClient()
	}

	resp, err := client.Req("HEAD", uri, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		pcsverbose.Verbosef("DEBUG: GetFileName ParseMediaType error: %s\n", err)
		return path.Base(uri), nil
	}

	filename, err = url.QueryUnescape(params["filename"])
	if err != nil {
		return
	}

	if filename == "" {
		filename = path.Base(uri)
	}

	return
}

// ParseContentRange 解析Content-Range
func ParseContentRange(contentRange string) (contentLength int64) {
	raw := ContentRangeRE.FindStringSubmatch(contentRange)
	if len(raw) < 2 {
		return -1
	}

	c, err := strconv.ParseInt(raw[1], 10, 64)
	if err != nil {
		return -1
	}
	return c
}

func fixCacheSize(size *int) {
	if *size < 1024 {
		*size = 1024
	}
}
