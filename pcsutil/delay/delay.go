package delay

import (
	"time"
)

// NewDelayChan 发送延时信号
func NewDelayChan(t time.Duration) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		time.Sleep(t)
		close(c)
	}()
	return c
}
