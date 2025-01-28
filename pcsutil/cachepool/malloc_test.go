package cachepool_test

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/cachepool"
	"runtime"
	"testing"
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
}
