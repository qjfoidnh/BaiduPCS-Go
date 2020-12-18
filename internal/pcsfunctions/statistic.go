package pcsfunctions

import (
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires"
	"sync/atomic"
	"time"
)

type (
	Statistic struct {
		totalSize int64
		startTime time.Time
	}
)

func (s *Statistic) AddTotalSize(size int64) int64 {
	return atomic.AddInt64(&s.totalSize, size)
}

func (s *Statistic) TotalSize() int64 {
	return s.totalSize
}

func (s *Statistic) StartTimer() {
	s.startTime = time.Now()
	expires.StripMono(&s.startTime)
}

func (s *Statistic) Elapsed() time.Duration {
	return time.Now().Sub(s.startTime)
}