// Package randominfo 提供随机信息生成服务
package randominfo

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// RandomNumber 返回[min, max]随机数字
func RandomNumber(min, max uint64) (v uint64) {
	if min > max {
		min, max = max, min
	}
	binary.Read(cryptorand.Reader, binary.BigEndian, &v)
	return v%(max-min) + min
}

// RandomBytes 随机字节数组
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	cryptorand.Read(b)
	return b
}

// RandomMD5String 随机md5字符串
func RandomMD5String() string {
	m := md5.New()
	m.Write(RamdomBytes(4))
	m.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 16)))
	return hex.EncodeToString(m.Sum(nil))
}

// RandomMD5UpperString 随机md5字符串, 大写
func RandomMD5UpperString() string {
	return strings.ToUpper(RamdomMD5String())
}

var (
	RamdomNumber         = RandomNumber
	RamdomBytes          = RandomBytes
	RamdomMD5String      = RandomMD5String
	RamdomMD5UpperString = RandomMD5String
)
