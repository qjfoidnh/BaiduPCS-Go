package uploader

import (
	"context"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/BaiduPCS-Go/requester/rio"
	"github.com/iikira/BaiduPCS-Go/requester/rio/speeds"
	"sync"
	"time"
)

type (
	// MultiUpload 支持多线程的上传, 可用于断点续传
	MultiUpload interface {
		Precreate() (perr error)
		TmpFile(ctx context.Context, partseq int, partOffset int64, readerlen64 rio.ReaderLen64) (checksum string, terr error)
		CreateSuperFile(checksumList ...string) (cerr error)
	}

	// MultiUploader 多线程上传
	MultiUploader struct {
		onExecuteEvent      requester.Event        //开始上传事件
		onSuccessEvent      requester.Event        //成功上传事件
		onFinishEvent       requester.Event        //结束上传事件
		onCancelEvent       requester.Event        //取消上传事件
		onErrorEvent        requester.EventOnError //上传出错事件
		onUploadStatusEvent UploadStatusFunc       //上传状态事件

		instanceState *InstanceState

		multiUpload MultiUpload       // 上传体接口
		file        rio.ReaderAtLen64 // 上传
		config      *MultiUploaderConfig
		workers     workerList
		speedsStat  *speeds.Speeds
		rateLimit   *speeds.RateLimit

		executeTime             time.Time
		finished                chan struct{}
		canceled                chan struct{}
		closeCanceledOnce       sync.Once
		updateInstanceStateChan chan struct{}
	}

	// MultiUploaderConfig 多线程上传配置
	MultiUploaderConfig struct {
		Parallel  int   // 上传并发量
		BlockSize int64 // 上传分块
		MaxRate   int64 // 限制最大上传速度
	}
)

// NewMultiUploader 初始化上传
func NewMultiUploader(multiUpload MultiUpload, file rio.ReaderAtLen64, config *MultiUploaderConfig) *MultiUploader {
	return &MultiUploader{
		multiUpload: multiUpload,
		file:        file,
		config:      config,
	}
}

// SetInstanceState 设置InstanceState, 断点续传信息
func (muer *MultiUploader) SetInstanceState(is *InstanceState) {
	muer.instanceState = is
}

func (muer *MultiUploader) lazyInit() {
	if muer.finished == nil {
		muer.finished = make(chan struct{}, 1)
	}
	if muer.canceled == nil {
		muer.canceled = make(chan struct{})
	}
	if muer.updateInstanceStateChan == nil {
		muer.updateInstanceStateChan = make(chan struct{}, 1)
	}
	if muer.config == nil {
		muer.config = &MultiUploaderConfig{}
	}
	if muer.config.Parallel <= 0 {
		muer.config.Parallel = 4
	}
	if muer.config.BlockSize <= 0 {
		muer.config.BlockSize = 1 * converter.GB
	}
	if muer.speedsStat == nil {
		muer.speedsStat = &speeds.Speeds{}
	}
}

func (muer *MultiUploader) check() {
	if muer.file == nil {
		panic("file is nil")
	}
	if muer.multiUpload == nil {
		panic("multiUpload is nil")
	}
}

// Execute 执行上传
func (muer *MultiUploader) Execute() {
	muer.check()
	muer.lazyInit()

	// 初始化限速
	if muer.config.MaxRate > 0 {
		muer.rateLimit = speeds.NewRateLimit(muer.config.MaxRate)
		defer muer.rateLimit.Stop()
	}

	// 分配任务
	if muer.instanceState != nil {
		muer.workers = muer.getWorkerListByInstanceState(muer.instanceState)
		uploaderVerbose.Infof("upload task CREATED from instance state\n")
	} else {
		muer.workers = muer.getWorkerListByInstanceState(&InstanceState{
			BlockList: SplitBlock(muer.file.Len(), muer.config.BlockSize),
		})

		uploaderVerbose.Infof("upload task CREATED: block size: %d, num: %d\n", muer.config.BlockSize, len(muer.workers))
	}

	// 开始上传
	muer.executeTime = time.Now()
	pcsutil.Trigger(muer.onExecuteEvent)

	muer.uploadStatusEvent()

	err := muer.upload()

	// 完成
	muer.finished <- struct{}{}
	if err != nil {
		if err == context.Canceled {
			if muer.onCancelEvent != nil {
				muer.onCancelEvent()
			}
		} else if muer.onErrorEvent != nil {
			muer.onErrorEvent(err)
		}
	} else {
		pcsutil.TriggerOnSync(muer.onSuccessEvent)
	}
	pcsutil.TriggerOnSync(muer.onFinishEvent)
}

// InstanceState 返回断点续传信息
func (muer *MultiUploader) InstanceState() *InstanceState {
	blockStates := make([]*BlockState, 0, len(muer.workers))
	for _, wer := range muer.workers {
		blockStates = append(blockStates, &BlockState{
			ID:       wer.id,
			Range:    wer.splitUnit.Range(),
			CheckSum: wer.checksum,
		})
	}
	return &InstanceState{
		BlockList: blockStates,
	}
}

// Cancel 取消上传
func (muer *MultiUploader) Cancel() {
	close(muer.canceled)
}

//OnExecute 设置开始上传事件
func (muer *MultiUploader) OnExecute(onExecuteEvent requester.Event) {
	muer.onExecuteEvent = onExecuteEvent
}

//OnSuccess 设置成功上传事件
func (muer *MultiUploader) OnSuccess(onSuccessEvent requester.Event) {
	muer.onSuccessEvent = onSuccessEvent
}

//OnFinish 设置结束上传事件
func (muer *MultiUploader) OnFinish(onFinishEvent requester.Event) {
	muer.onFinishEvent = onFinishEvent
}

//OnCancel 设置取消上传事件
func (muer *MultiUploader) OnCancel(onCancelEvent requester.Event) {
	muer.onCancelEvent = onCancelEvent
}

//OnError 设置上传发生错误事件
func (muer *MultiUploader) OnError(onErrorEvent requester.EventOnError) {
	muer.onErrorEvent = onErrorEvent
}

//OnUploadStatusEvent 设置上传状态事件
func (muer *MultiUploader) OnUploadStatusEvent(f UploadStatusFunc) {
	muer.onUploadStatusEvent = f
}
