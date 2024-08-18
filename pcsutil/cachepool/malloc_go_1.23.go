//go:build go1.23
// +build go1.23

package cachepool

import (
	"unsafe"
)

// 说明：
// 由于GO 1.23版本取消了 go:linkname 的支持，所以1.23以及以上版本需要使用本文件替代原始文件 malloc.go

// RawByteSlice point to runtime.rawbyteslice
func RawByteSlice(size int) (b []byte) {
	bytesArray := make([]byte, size)
	return bytesArray
}

// RawMalloc allocates a new slice. The slice is not zeroed.
func RawMalloc(size int) unsafe.Pointer {
	bytesArray := make([]byte, size)
	// 使用unsafe.Pointer获取字节数组的指针
	bytesPtr := unsafe.Pointer(&bytesArray[0])
	return bytesPtr
}

// RawMallocByteSlice allocates a new byte slice. The slice is not zeroed.
func RawMallocByteSlice(size int) []byte {
	bytesArray := make([]byte, size)
	return bytesArray
}