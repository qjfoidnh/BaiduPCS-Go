package uploader

import (
	"github.com/iikira/BaiduPCS-Go/requester/rio"
	"sync/atomic"
)

type (
	// Readed64 增加获取已读取数据量, 用于统计速度
	Readed64 interface {
		rio.ReaderLen64
		Readed() int64
	}

	readed64 struct {
		readed int64
		rio.ReaderLen64
	}
)

// NewReaded64 实现Readed64接口
func NewReaded64(rl rio.ReaderLen64) Readed64 {
	return &readed64{
		readed:      0,
		ReaderLen64: rl,
	}
}

func (r64 *readed64) Read(p []byte) (n int, err error) {
	n, err = r64.ReaderLen64.Read(p)
	atomic.AddInt64(&r64.readed, int64(n))
	return n, err
}

func (r64 *readed64) Readed() int64 {
	return atomic.LoadInt64(&r64.readed)
}
