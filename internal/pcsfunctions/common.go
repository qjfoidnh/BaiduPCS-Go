package pcsfunctions

import (
	"time"
)

// RetryWait 失败重试等待事件
func RetryWait(retry int) time.Duration {
	if retry < 3 {
		return 2 * time.Duration(retry) * time.Second
	}
	return 6 * time.Second
}
