package cachemap_test

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires/cachemap"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheMapDataExpires(t *testing.T) {
	cm := cachemap.CacheOpMap{}
	cache := cm.LazyInitCachePoolOp("op")
	cache.Store("key_1", expires.NewDataExpires("value_1", 1*time.Second))

	time.Sleep(2 * time.Second)
	data, ok := cache.Load("key_1")
	if ok {
		fmt.Printf("data: %s\n", data.Data())
		// 超时仍能读取到数据, 失败
		t.FailNow()
	}
}

func TestCacheOperation(t *testing.T) {
	cm := cachemap.CacheOpMap{}
	data := cm.CacheOperation("op", "key_1", func() expires.DataExpires {
		return expires.NewDataExpires("value_1", 1*time.Second)
	})
	fmt.Printf("data: %s\n", data.Data())

	newData := cm.CacheOperation("op", "key_1", func() expires.DataExpires {
		return expires.NewDataExpires("value_3", 1*time.Second)
	})
	if data != newData {
		t.FailNow()
	}
	fmt.Printf("data: %s\n", data.Data())
}

func TestCacheOperation_LockKey(t *testing.T) {
	cm := cachemap.CacheOpMap{}
	wg := sync.WaitGroup{}
	wg.Add(5000)

	var (
		execTimes1 int32 = 0 // 执行次数1
		execTimes2 int32 = 0 // 执行次数2
	)

	for i := 0; i < 5000; i++ {
		go func(i int) {
			defer wg.Done()
			cm.CacheOperation("op", "key_1", func() expires.DataExpires {
				time.Sleep(50 * time.Microsecond) // 一些耗时的操作
				atomic.AddInt32(&execTimes1, 1)
				return expires.NewDataExpires(fmt.Sprintf("value_1: %d", i), 10*time.Second)
			})

			cm.CacheOperation("op", "key_2", func() expires.DataExpires {
				time.Sleep(50 * time.Microsecond) // 一些耗时的操作
				atomic.AddInt32(&execTimes2, 1)
				return expires.NewDataExpires(fmt.Sprintf("value_2: %d", i), 10*time.Second)
			})
		}(i)
	}
	wg.Wait()

	// 执行次数应为1
	if execTimes1 != 1 {
		fmt.Printf("execTimes1: %d\n", execTimes1)
		t.FailNow()
	}
	if execTimes2 != 1 {
		fmt.Printf("execTimes2: %d\n", execTimes2)
		t.FailNow()
	}
}
