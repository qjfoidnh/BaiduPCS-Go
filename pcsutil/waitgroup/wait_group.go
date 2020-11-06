// Package waitgroup sync.WaitGroup extension
package waitgroup

import "sync"

// WaitGroup 在 sync.WaitGroup 的基础上, 新增线程控制功能
type WaitGroup struct {
	wg sync.WaitGroup
	p  chan struct{}

	sync.RWMutex
}

// NewWaitGroup returns a pointer to a new `WaitGroup` object.
// parallel 为最大并发数, 0 代表无限制
func NewWaitGroup(parallel int) (w *WaitGroup) {
	w = &WaitGroup{
		wg: sync.WaitGroup{},
	}

	if parallel <= 0 {
		return
	}

	w.p = make(chan struct{}, parallel)
	return
}

// AddDelta sync.WaitGroup.Add(1)
func (w *WaitGroup) AddDelta() {
	if w.p != nil {
		w.p <- struct{}{}
	}

	w.wg.Add(1)
}

// Done sync.WaitGroup.Done()
func (w *WaitGroup) Done() {
	w.wg.Done()

	if w.p != nil {
		<-w.p
	}
}

// Wait 参照 sync.WaitGroup 的 Wait 方法
func (w *WaitGroup) Wait() {
	w.wg.Wait()
	if w.p != nil {
		close(w.p)
	}
}

// Parallel 返回当前正在进行的任务数量
func (w *WaitGroup) Parallel() int {
	return len(w.p)
}
