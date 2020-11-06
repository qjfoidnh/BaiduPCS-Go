package netdisksign

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/iikira/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
)

func DevUID(feature string) string {
	m := md5.New()
	m.Write(converter.ToBytes(feature))
	res := m.Sum(nil)
	resHex := cachepool.RawMallocByteSlice(34)
	hex.Encode(resHex[2:], res)
	resHex[0] = 'O'
	resHex[1] = '|'
	return converter.ToString(bytes.ToUpper(resHex))
}
