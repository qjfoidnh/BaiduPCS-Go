package cachepool

import (
	"reflect"
	"unsafe"
)

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
