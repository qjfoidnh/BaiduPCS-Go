package cachemap

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/expires"
	"sync"
)

type (
	CacheUnit interface {
		Delete(key interface{})
		Load(key interface{}) (value expires.DataExpires, ok bool)
		LoadOrStore(key interface{}, value expires.DataExpires) (actual expires.DataExpires, loaded bool)
		Range(f func(key interface{}, value expires.DataExpires) bool)
		Store(key interface{}, value expires.DataExpires)
		LockKey(key interface{})
		UnlockKey(key interface{})
	}

	cacheUnit struct {
		unit   sync.Map
		keyMap sync.Map
	}
)

func (cu *cacheUnit) Delete(key interface{}) {
	cu.unit.Delete(key)
	cu.keyMap.Delete(key)
}

func (cu *cacheUnit) Load(key interface{}) (value expires.DataExpires, ok bool) {
	val, ok := cu.unit.Load(key)
	if !ok {
		return nil, ok
	}
	exp := val.(expires.DataExpires)
	if exp.IsExpires() {
		cu.unit.Delete(key)
		return nil, false
	}
	return exp, ok
}

func (cu *cacheUnit) Range(f func(key interface{}, value expires.DataExpires) bool) {
	cu.unit.Range(func(k, val interface{}) bool {
		exp := val.(expires.DataExpires)
		if exp.IsExpires() {
			cu.unit.Delete(k)
			return true
		}
		return f(k, val.(expires.DataExpires))
	})
}

func (cu *cacheUnit) LoadOrStore(key interface{}, value expires.DataExpires) (actual expires.DataExpires, loaded bool) {
	ac, loaded := cu.unit.LoadOrStore(key, value)
	exp := ac.(expires.DataExpires)
	if exp.IsExpires() {
		cu.unit.Delete(key)
		return nil, false
	}
	return exp, loaded
}

func (cu *cacheUnit) Store(key interface{}, value expires.DataExpires) {
	if value.IsExpires() {
		return
	}
	cu.unit.Store(key, value)
}

func (cu *cacheUnit) LockKey(key interface{}) {
	muItf, _ := cu.keyMap.LoadOrStore(key, &sync.Mutex{})
	mu := muItf.(*sync.Mutex)
	mu.Lock()
}

func (cu *cacheUnit) UnlockKey(key interface{}) {
	muItf, _ := cu.keyMap.LoadOrStore(key, &sync.Mutex{})
	mu := muItf.(*sync.Mutex)
	mu.Unlock()
}
