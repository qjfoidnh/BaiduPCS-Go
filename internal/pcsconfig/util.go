package pcsconfig

import (
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"strings"
)

// AverageParallel 返回平均的下载最大并发量
func AverageParallel(parallel, downloadLoad int) int {
	if downloadLoad < 1 {
		return 1
	}

	p := parallel / downloadLoad
	if p < 1 {
		return 1
	}
	return p
}

func stripPerSecond(sizeStr string) string {
	i := strings.LastIndex(sizeStr, "/")
	if i < 0 {
		return sizeStr
	}
	return sizeStr[:i]
}

func showMaxRate(size int64) string {
	if size <= 0 {
		return "不限制"
	}
	return converter.ConvertFileSize(size, 2) + "/s"
}
