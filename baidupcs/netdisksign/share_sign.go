package netdisksign

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"strconv"
)

func ShareSURLInfoSign(shareID int64) []byte {
	s := strconv.FormatInt(shareID, 10)
	m := md5.New()
	m.Write(converter.ToBytes(s))
	m.Write([]byte("_sharesurlinfo!@#"))
	res := m.Sum(nil)
	resHex := cachepool.RawMallocByteSlice(32)
	hex.Encode(resHex, res)
	return resHex
}
