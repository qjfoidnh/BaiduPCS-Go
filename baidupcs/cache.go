package baidupcs

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs/expires"
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
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
	data := pcs.cacheOpMap.CacheOperation(OperationFilesDirectoriesList, path+"_"+fmt.Sprint(options), func() expires.DataExpires {
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
