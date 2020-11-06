package rio

import (
	"unsafe"
)

// Buffer 为固定长度的 Buf, 实现 io.WriterAt 接口
type Buffer struct {
	Buf []byte
}

// NewBuffer 初始化 Buffer
func NewBuffer(buf []byte) *Buffer {
	return &Buffer{
		Buf: buf,
	}
}

// ReadAt 实现 io.ReadAt 接口
// 不进行越界检查
func (b *Buffer) ReadAt(p []byte, off int64) (n int, err error) {
	n = copy(p, b.Buf[off:])
	return n, nil
}

// WriteAt 实现 io.WriterAt 接口
// 不进行越界检查
func (b *Buffer) WriteAt(p []byte, off int64) (n int, err error) {
	n = copy(b.Buf[off:], p)
	return n, nil
}

// Bytes 返回 buf
func (b *Buffer) Bytes() []byte {
	return b.Buf
}

func (b *Buffer) String() string {
	return *(*string)(unsafe.Pointer(&b.Buf))
}
