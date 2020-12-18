package transfer

import (
	"github.com/qjfoidnh/BaiduPCS-Go/requester/rio/speeds"
	"sync"
	"sync/atomic"
	"time"
)

type (
	//DownloadStatuser 下载状态接口
	DownloadStatuser interface {
		TotalSize() int64
		Downloaded() int64
		SpeedsPerSecond() int64
		TimeElapsed() time.Duration // 已开始时间
		TimeLeft() time.Duration    // 预计剩余时间, 负数代表未知
	}

	//DownloadStatus 下载状态及统计信息
	DownloadStatus struct {
		totalSize        int64         // 总大小
		downloaded       int64         // 已下载的数据量
		speedsDownloaded int64         // 用于统计速度的downloaded
		maxSpeeds        int64         // 最大下载速度
		tmpSpeeds        int64         // 缓存的速度
		speedsStat       speeds.Speeds // 速度统计 (注意对齐)

		startTime time.Time // 开始下载的时间

		rateLimit *speeds.RateLimit // 限速控制

		gen *RangeListGen // Range生成状态
		mu  sync.Mutex
	}
)

//NewDownloadStatus 初始化DownloadStatus
func NewDownloadStatus() *DownloadStatus {
	return &DownloadStatus{
		startTime: time.Now(),
	}
}

// SetRateLimit 设置限速
func (ds *DownloadStatus) SetRateLimit(rl *speeds.RateLimit) {
	ds.rateLimit = rl
}

//SetTotalSize 返回总大小
func (ds *DownloadStatus) SetTotalSize(size int64) {
	ds.totalSize = size
}

//AddDownloaded 增加已下载数据量
func (ds *DownloadStatus) AddDownloaded(d int64) {
	atomic.AddInt64(&ds.downloaded, d)
}

//AddTotalSize 增加总大小 (不支持多线程)
func (ds *DownloadStatus) AddTotalSize(size int64) {
	ds.totalSize += size
}

//AddSpeedsDownloaded 增加已下载数据量, 用于统计速度
func (ds *DownloadStatus) AddSpeedsDownloaded(d int64) {
	if ds.rateLimit != nil {
		ds.rateLimit.Add(d)
	}
	ds.speedsStat.Add(d)
}

//SetMaxSpeeds 设置最大速度, 原子操作
func (ds *DownloadStatus) SetMaxSpeeds(speeds int64) {
	if speeds > atomic.LoadInt64(&ds.maxSpeeds) {
		atomic.StoreInt64(&ds.maxSpeeds, speeds)
	}
}

//ClearMaxSpeeds 清空统计最大速度, 原子操作
func (ds *DownloadStatus) ClearMaxSpeeds() {
	atomic.StoreInt64(&ds.maxSpeeds, 0)
}

//TotalSize 返回总大小
func (ds *DownloadStatus) TotalSize() int64 {
	return ds.totalSize
}

//Downloaded 返回已下载数据量
func (ds *DownloadStatus) Downloaded() int64 {
	return atomic.LoadInt64(&ds.downloaded)
}

// UpdateSpeeds 更新speeds
func (ds *DownloadStatus) UpdateSpeeds() {
	atomic.StoreInt64(&ds.tmpSpeeds, ds.speedsStat.GetSpeeds())
}

//SpeedsPerSecond 返回每秒速度
func (ds *DownloadStatus) SpeedsPerSecond() int64 {
	return atomic.LoadInt64(&ds.tmpSpeeds)
}

//MaxSpeeds 返回最大速度
func (ds *DownloadStatus) MaxSpeeds() int64 {
	return atomic.LoadInt64(&ds.maxSpeeds)
}

//TimeElapsed 返回花费的时间
func (ds *DownloadStatus) TimeElapsed() (elapsed time.Duration) {
	return time.Since(ds.startTime)
}

//TimeLeft 返回预计剩余时间
func (ds *DownloadStatus) TimeLeft() (left time.Duration) {
	speeds := atomic.LoadInt64(&ds.tmpSpeeds)
	if speeds <= 0 {
		left = -1
	} else {
		left = time.Duration((ds.totalSize-ds.downloaded)/(speeds)) * time.Second
	}
	return
}

// RangeListGen 返回RangeListGen
func (ds *DownloadStatus) RangeListGen() *RangeListGen {
	return ds.gen
}

// SetRangeListGen 设置RangeListGen
func (ds *DownloadStatus) SetRangeListGen(gen *RangeListGen) {
	ds.gen = gen
}
