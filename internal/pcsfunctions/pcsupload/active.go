package pcsupload

import (
	"github.com/qjfoidnh/BaiduPCS-Go/requester/uploader"
	"sync"
	"sync/atomic"
)

var (
	// activeUploaders 记录当前正在执行的上传器
	activeUploaders sync.Map
	// uploadStopping 全局优雅停止标志
	uploadStopping int32
)

// GracefulStopActiveUploaders 优雅停止所有进行中的上传
// 不再派发新分片, 进行中的分片继续执行直至完成
func GracefulStopActiveUploaders() {
	atomic.StoreInt32(&uploadStopping, 1)
	activeUploaders.Range(func(k, _ interface{}) bool {
		k.(*uploader.MultiUploader).GracefulStop()
		return true
	})
}

// RegisterActiveUploader 注册进行中的上传器
// 若注册时已处于停止状态, 立即对该上传器生效, 避免"先停止后注册"的竞态
func RegisterActiveUploader(muer *uploader.MultiUploader) {
	activeUploaders.Store(muer, struct{}{})
	if atomic.LoadInt32(&uploadStopping) == 1 {
		muer.GracefulStop()
	}
}

// UnregisterActiveUploader 注销上传器
func UnregisterActiveUploader(muer *uploader.MultiUploader) {
	activeUploaders.Delete(muer)
}
