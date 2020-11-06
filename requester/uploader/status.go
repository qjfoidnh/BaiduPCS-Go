package uploader

import (
	"time"
)

type (
	// Status 上传状态接口
	Status interface {
		TotalSize() int64           // 总大小
		Uploaded() int64            // 已上传数据
		SpeedsPerSecond() int64     // 每秒的上传速度
		TimeElapsed() time.Duration // 上传时间
	}

	// UploadStatus 上传状态
	UploadStatus struct {
		totalSize       int64         // 总大小
		uploaded        int64         // 已上传数据
		speedsPerSecond int64         // 每秒的上传速度
		timeElapsed     time.Duration // 上传时间
	}

	UploadStatusFunc func(status Status, updateChan <-chan struct{})
)

// TotalSize 返回总大小
func (us *UploadStatus) TotalSize() int64 {
	return us.totalSize
}

// Uploaded 返回已上传数据
func (us *UploadStatus) Uploaded() int64 {
	return us.uploaded
}

// SpeedsPerSecond 返回每秒的上传速度
func (us *UploadStatus) SpeedsPerSecond() int64 {
	return us.speedsPerSecond
}

// TimeElapsed 返回上传时间
func (us *UploadStatus) TimeElapsed() time.Duration {
	return us.timeElapsed
}

// GetStatusChan 获取上传状态
func (u *Uploader) GetStatusChan() <-chan Status {
	c := make(chan Status)

	go func() {
		for {
			select {
			case <-u.finished:
				close(c)
				return
			default:
				if !u.executed {
					time.Sleep(1 * time.Second)
					continue
				}

				old := u.readed64.Readed()
				time.Sleep(1 * time.Second) // 每秒统计

				readed := u.readed64.Readed()
				c <- &UploadStatus{
					totalSize:       u.readed64.Len(),
					uploaded:        readed,
					speedsPerSecond: readed - old,
					timeElapsed:     time.Since(u.executeTime) / 1e7 * 1e7,
				}
			}
		}
	}()
	return c
}

func (muer *MultiUploader) uploadStatusEvent() {
	if muer.onUploadStatusEvent == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second) // 每秒统计
		defer ticker.Stop()
		for {
			select {
			case <-muer.finished:
				return
			case <-ticker.C:
				readed := muer.workers.Readed()
				muer.onUploadStatusEvent(&UploadStatus{
					totalSize:       muer.file.Len(),
					uploaded:        readed,
					speedsPerSecond: muer.speedsStat.GetSpeeds(),
					timeElapsed:     time.Since(muer.executeTime) / 1e8 * 1e8,
				}, muer.updateInstanceStateChan)
			}
		}
	}()
}
