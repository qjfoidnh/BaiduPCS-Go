//go:build !go1.23
// +build !go1.23

package cachepool

import (
	"reflect"
	"unsafe"
)

// 函数声明可以省略主体。 这样的声明为Go外部实现的功能（例如汇编例程）提供了签名。这是在汇编中实现函数的方式。

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

//go:linkname rawbyteslice runtime.rawbyteslice
func rawbyteslice(size int) (b []byte)

// RawByteSlice point to runtime.rawbyteslice
func RawByteSlice(size int) (b []byte) {
	return rawbyteslice(size)
}

// RawMalloc allocates a new slice. The slice is not zeroed.
func RawMalloc(size int) unsafe.Pointer {
	return mallocgc(uintptr(size), 0, false)
}

// RawMallocByteSlice allocates a new byte slice. The slice is not zeroed.
func RawMallocByteSlice(size int) []byte {
	p := mallocgc(uintptr(size), 0, false)
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p),
		Len:  size,
		Cap:  size,
	}))
	return b
}