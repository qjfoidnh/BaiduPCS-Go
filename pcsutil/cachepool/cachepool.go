package cachepool

import (
	"sync"
)

var (
	//CachePool []byte 缓存池 2
	CachePool = cachePool2{}
)

//Cache cache
type Cache interface {
	Bytes() []byte
	Free()
}

type cache struct {
	isUsed bool
	b      []byte
}

func (c *cache) Bytes() []byte {
	if !c.isUsed {
		return nil
	}
	return c.b
}

func (c *cache) Free() {
	c.isUsed = false
}

type cachePool2 struct {
	pool []*cache
	mu   sync.Mutex
}

func (cp2 *cachePool2) Require(size int) Cache {
	cp2.mu.Lock()
	defer cp2.mu.Unlock()
	for k := range cp2.pool {
		if cp2.pool[k] == nil || cp2.pool[k].isUsed || len(cp2.pool[k].b) < size {
			continue
		}

		cp2.pool[k].isUsed = true
		return cp2.pool[k]
	}
	newCache := &cache{
		isUsed: true,
		b:      RawMallocByteSlice(size),
	}
	cp2.addCache(newCache)
	return newCache
}

func (cp2 *cachePool2) addCache(newCache *cache) {
	for k := range cp2.pool {
		if cp2.pool[k] == nil {
			cp2.pool[k] = newCache
			return
		}
	}
	cp2.pool = append(cp2.pool, newCache)
}

func (cp2 *cachePool2) DeleteNotUsed() {
	cp2.mu.Lock()
	defer cp2.mu.Unlock()
	for k := range cp2.pool {
		if cp2.pool[k] == nil {
			continue
		}

		if !cp2.pool[k].isUsed {
			cp2.pool[k] = nil
		}
	}
}

func (cp2 *cachePool2) DeleteAll() {
	cp2.mu.Lock()
	defer cp2.mu.Unlock()
	for k := range cp2.pool {
		cp2.pool[k] = nil
	}
}

//Require 申请Cache
func Require(size int) Cache {
	return CachePool.Require(size)
}
