package cachemap

import (
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires"
)

type (
	OpFunc          func() expires.DataExpires
	OpFuncWithError func() (expires.DataExpires, error)
)

func (cm *CacheOpMap) CacheOperation(op string, key interface{}, opFunc OpFunc) (data expires.DataExpires) {
	var (
		cache = cm.LazyInitCachePoolOp(op)
		ok    bool
	)

	cache.LockKey(key)
	defer cache.UnlockKey(key)
	data, ok = cache.Load(key)
	if !ok {
		data = opFunc()
		if data != nil {
			cache.Store(key, data)
		}
		return
	}

	return
}

func (cm *CacheOpMap) CacheOperationWithError(op string, key interface{}, opFunc OpFuncWithError) (data expires.DataExpires, err error) {
	var (
		cache = cm.LazyInitCachePoolOp(op)
		ok    bool
	)

	cache.LockKey(key)
	defer cache.UnlockKey(key)
	data, ok = cache.Load(key)
	if !ok {
		data, err = opFunc()
		if err != nil {
			return
		}
		if data == nil {
			// 数据为空时也不存
			return
		}
		cache.Store(key, data)
	}

	return
}
