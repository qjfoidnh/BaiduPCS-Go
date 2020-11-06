// Package speeds 速度计算工具包
package speeds

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	// Speeds 统计速度
	Speeds struct {
		count    int64
		interval time.Duration // 刷新周期
		nowTime  time.Time
		once     sync.Once
	}
)

func (sps *Speeds) initOnce() {
	sps.once.Do(func() {
		sps.nowTime = time.Now()
		if sps.interval <= 0 {
			sps.interval = 1 * time.Second
		}
	})
}

// SetInterval 设置刷新周期
func (sps *Speeds) SetInterval(interval time.Duration) {
	if interval <= 0 {
		return
	}
	sps.interval = interval
}

// Add 原子操作, 增加数据量
func (sps *Speeds) Add(count int64) {
	// 初始化
	sps.initOnce()
	atomic.AddInt64(&sps.count, count)
}

// GetSpeeds 结束统计速度, 并返回速度
func (sps *Speeds) GetSpeeds() (speeds int64) {
	sps.initOnce()

	since := time.Since(sps.nowTime)
	if since <= 0 {
		return 0
	}
	speeds = int64(float64(atomic.LoadInt64(&sps.count)) * sps.interval.Seconds() / since.Seconds())

	// 更新下一轮
	if since >= sps.interval {
		atomic.StoreInt64(&sps.count, 0)
		sps.nowTime = time.Now()
	}
	return
}
