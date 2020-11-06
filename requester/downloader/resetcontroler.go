package downloader

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/expires"
	"sync"
	"time"
)

// ResetController 网络连接控制器
type ResetController struct {
	mu          sync.Mutex
	currentTime time.Time
	maxResetNum int
	resetEntity map[expires.Expires]struct{}
}

// NewResetController 初始化*ResetController
func NewResetController(maxResetNum int) *ResetController {
	return &ResetController{
		currentTime: time.Now(),
		maxResetNum: maxResetNum,
		resetEntity: map[expires.Expires]struct{}{},
	}
}

func (rc *ResetController) update() {
	for k := range rc.resetEntity {
		if k.IsExpires() {
			delete(rc.resetEntity, k)
		}
	}
}

// AddResetNum 增加连接
func (rc *ResetController) AddResetNum() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.update()
	rc.resetEntity[expires.NewExpires(9*time.Second)] = struct{}{}
}

// CanReset 是否可以建立连接
func (rc *ResetController) CanReset() bool {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.update()
	return len(rc.resetEntity) < rc.maxResetNum
}
