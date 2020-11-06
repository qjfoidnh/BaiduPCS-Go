package speeds

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	RateLimit struct {
		MaxRate int64

		count           int64
		interval        time.Duration
		ticker          *time.Ticker
		muChan          chan struct{}
		closeChan       chan struct{}
		backServiceOnce sync.Once
	}

	// AddCountFunc func() (count int64)
)

func NewRateLimit(maxRate int64) *RateLimit {
	return &RateLimit{
		MaxRate: maxRate,
	}
}

func (rl *RateLimit) SetInterval(i time.Duration) {
	if i <= 0 {
		i = 1 * time.Second
	}
	rl.interval = i
	if rl.ticker != nil {
		rl.ticker.Stop()
		rl.ticker = time.NewTicker(i)
	}
}

func (rl *RateLimit) Stop() {
	if rl.ticker != nil {
		rl.ticker.Stop()
	}
	if rl.closeChan != nil {
		close(rl.closeChan)
	}
	return
}

func (rl *RateLimit) resetChan() {
	if rl.muChan != nil {
		close(rl.muChan)
	}
	rl.muChan = make(chan struct{})
}

func (rl *RateLimit) backService() {
	if rl.interval <= 0 {
		rl.interval = 1 * time.Second
	}
	rl.ticker = time.NewTicker(rl.interval)
	rl.closeChan = make(chan struct{})
	rl.resetChan()
	go func() {
		for {
			select {
			case <-rl.ticker.C:
				rl.resetChan()
				atomic.StoreInt64(&rl.count, 0)
			case <-rl.closeChan:
				return
			}
		}
	}()
}

func (rl *RateLimit) Add(count int64) {
	rl.backServiceOnce.Do(rl.backService)
	for {
		if atomic.LoadInt64(&rl.count) >= rl.MaxRate { // 超出最大限额
			// 阻塞
			<-rl.muChan
			continue
		}
		atomic.AddInt64(&rl.count, count)
		break
	}
}
