package baidupcs

import (
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"time"
)

// deleteCache 删除含有 dirs 的缓存
func (pcs *BaiduPCS) deleteCache(dirs []string) {
	cache := pcs.cacheOpMap.LazyInitCachePoolOp(OperationFilesDirectoriesList)
	for _, v := range dirs {
		key := v + "_" + defaultOrderOptionsStr
		_, ok := cache.Load(key)
		if ok {
			cache.Delete(key)
		}
	}
}

// CacheFilesDirectoriesList 缓存获取
func (pcs *BaiduPCS) CacheFilesDirectoriesList(path string, options *OrderOptions) (fdl FileDirectoryList, pcsError pcserror.Error) {
	data := pcs.cacheOpMap.CacheOperation(OperationFilesDirectoriesList, path+"_"+string(options.By)+string(options.Order), func() expires.DataExpires {
		fdl, pcsError = pcs.FilesDirectoriesList(path, options)
		if pcsError != nil {
			return nil
		}
		return expires.NewDataExpires(fdl, 1*time.Minute)
	})
	if pcsError != nil {
		return
	}
	return data.Data().(FileDirectoryList), nil
}

// CacheUK 缓存获取
func (pcs *BaiduPCS) CacheUK() (uk int64, pcsError pcserror.Error) {
	data := pcs.cacheOpMap.CacheOperation(OperationGetUK, pcs.GetBDUSS(), func() expires.DataExpires {
		uk, pcsError = pcs.UK()
		if pcsError != nil {
			return nil
		}
		return expires.NewDataExpires(uk, 24*time.Hour)
	})
	if pcsError != nil {
		return
	}
	return data.Data().(int64), nil
}
