package pcsupload

import (
	"crypto/md5"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"io"
	"strconv"
)

func getBlockSize(fileSize int64) int64 {
	blockSize := baidupcs.MinUploadBlockSize
	if fileSize > 1*converter.GB && fileSize < 4*converter.GB {
		blockSize = 2 * baidupcs.MinUploadBlockSize
	} else if fileSize >= 4*converter.GB {
		blockSize = baidupcs.MaxUploadBlockSize
	}
	return blockSize
}

func creaetDataOffset(contentMD5 string, uk, dataTime, fileSize, subSize int64) (offset int64, err error) {
	h := md5.New()
	ts := strconv.FormatInt(dataTime, 10)
	sumStr := fmt.Sprintf("%d%s%s", uk, contentMD5, ts)
	io.WriteString(h, sumStr)
	mixedMD5 := fmt.Sprintf("%x", h.Sum(nil))
	rawOffset, err := strconv.ParseInt(mixedMD5[0:8], 16, 64)
	if err != nil {
		return
	}
	if fileSize-subSize+1 <= 1 {
		offset = 0
		return
	}
	offset = rawOffset % (fileSize - subSize + 1)
	return
}
