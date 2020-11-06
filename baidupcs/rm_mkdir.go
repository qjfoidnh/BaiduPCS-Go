package baidupcs

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
	"path"
)

// Remove 批量删除文件/目录
func (pcs *BaiduPCS) Remove(paths ...string) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareRemove(paths...)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.DecodePCSJSONError(OperationRemove, dataReadCloser)
	if errInfo != nil {
		return errInfo
	}

	// 更新缓存
	pcs.deleteCache(allRelatedDir(paths))
	return nil
}

// Mkdir 创建目录
func (pcs *BaiduPCS) Mkdir(pcspath string) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareMkdir(pcspath)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.DecodePCSJSONError(OperationMkdir, dataReadCloser)
	if errInfo != nil {
		return errInfo
	}

	// 更新缓存
	pcs.deleteCache([]string{path.Dir(pcspath)})
	return
}
