package rio

import (
	cryptorand "crypto/rand"
	"io"
	"os"
	"sync/atomic"
)

type (
	fileReadedlen64 struct {
		readed int64
		f      *os.File
	}

	rdReadedlen64 struct {
		readed int64
		size   int64
		rd     io.Reader
	}
)

// NewFileReaderLen64 *os.File 实现 ReadedLen64 接口
func NewFileReaderLen64(f *os.File) ReaderLen64 {
	if f == nil {
		return nil
	}

	return &fileReadedlen64{
		f: f,
	}
}

// NewFileReaderAtLen64 *os.File 实现 ReaderAtLen64 接口
func NewFileReaderAtLen64(f *os.File) ReaderAtLen64 {
	if f == nil {
		return nil
	}

	return &fileReadedlen64{
		f: f,
	}
}

func NewCryptoRandReaderAtLen64(size int64) ReaderAtLen64 {
	return &rdReadedlen64{
		rd:   cryptorand.Reader,
		size: size,
	}
}

// Read 读文件, 并记录已读取数据量
func (fr *fileReadedlen64) Read(b []byte) (n int, err error) {
	n, err = fr.f.Read(b)
	atomic.AddInt64(&fr.readed, int64(n))
	return n, err
}

// ReadAt 读文件, 不记录已读取数据量
func (fr *fileReadedlen64) ReadAt(b []byte, off int64) (n int, err error) {
	n, err = fr.f.ReadAt(b, off)
	return n, err
}

// Len 返回文件的大小
func (fr *fileReadedlen64) Len() int64 {
	info, err := fr.f.Stat()
	if err != nil {
		return 0
	}
	return info.Size() - fr.readed
}

func (rr *rdReadedlen64) Read(b []byte) (n int, err error) {
	n, err = rr.ReadAt(b, 0)
	atomic.AddInt64(&rr.readed, int64(n))
	return n, err
}

func (rr *rdReadedlen64) ReadAt(b []byte, off int64) (n int, err error) {
	n, err = rr.rd.Read(b)
	return n, err
}

func (rr *rdReadedlen64) Len() int64 {
	return rr.size - rr.readed
}
