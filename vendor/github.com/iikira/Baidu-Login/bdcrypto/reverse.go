package bdcrypto

import (
	"unsafe"
)

// BytesReverse 反转字节数组, 此操作会修改原值
func BytesReverse(b []byte) []byte {
	length := len(b)
	for i := 0; i < length/2; i++ {
		b[i], b[length-i-1] = b[length-i-1], b[i]
	}
	return b
}

// StringReverse 反转字符串, 此操作不会修改原值
func StringReverse(s string) string {
	newBytes := make([]byte, len(s))
	copy(newBytes, s)
	b := BytesReverse(newBytes)
	return *(*string)(unsafe.Pointer(&b))
}
