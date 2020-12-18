package cachepool_test

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/cachepool"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

func TestMalloc(t *testing.T) {
	b := cachepool.RawMallocByteSlice(128)
	for k := range b {
		b[k] = byte(k)
	}
	fmt.Println(b)
	runtime.GC()

	b = cachepool.RawMallocByteSlice(128)
	fmt.Printf("---%s---\n", b)
	runtime.GC()

	b = cachepool.RawByteSlice(128)
	fmt.Println(b)
	runtime.GC()

	b = cachepool.RawByteSlice(127)
	bH := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	fmt.Printf("%#v\n", bH)
}
