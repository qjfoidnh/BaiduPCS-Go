package uploader

import (
	"bufio"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/rio/speeds"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/transfer"
	"io"
	"os"
	"sync"
)

type (
	// SplitUnit 将 io.ReaderAt 分割单元
	SplitUnit interface {
		Readed64
		io.Seeker
		Range() transfer.Range
		Left() int64
	}

	fileBlock struct {
		readRange     transfer.Range
		readed        int64
		readerAt      io.ReaderAt
		speedsStatRef *speeds.Speeds
		rateLimit     *speeds.RateLimit
		mu            sync.Mutex
	}

	bufioFileBlock struct {
		*fileBlock
		bufio *bufio.Reader
	}
)

// SplitBlock 文件分块
func SplitBlock(fileSize, blockSize int64) (blockList []*BlockState) {
	gen := transfer.NewRangeListGenBlockSize(fileSize, 0, blockSize)
	rangeCount := gen.RangeCount()
	blockList = make([]*BlockState, 0, rangeCount)
	for i := 0; i < rangeCount; i++ {
		id, r := gen.GenRange()
		blockList = append(blockList, &BlockState{
			ID:    id,
			Range: *r,
		})
	}
	return
}

// NewBufioSplitUnit io.ReaderAt实现SplitUnit接口, 有Buffer支持
func NewBufioSplitUnit(readerAt io.ReaderAt, readRange transfer.Range, speedsStat *speeds.Speeds, rateLimit *speeds.RateLimit) SplitUnit {
	su := &fileBlock{
		readerAt:      readerAt,
		readRange:     readRange,
		speedsStatRef: speedsStat,
		rateLimit:     rateLimit,
	}
	return &bufioFileBlock{
		fileBlock: su,
		bufio:     bufio.NewReaderSize(su, BufioReadSize),
	}
}

func (bfb *bufioFileBlock) Read(b []byte) (n int, err error) {
	return bfb.bufio.Read(b) // 间接调用fileBlock 的Read
}

// Read 只允许一个线程读同一个文件
func (fb *fileBlock) Read(b []byte) (n int, err error) {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	left := int(fb.Left())
	if left <= 0 {
		return 0, io.EOF
	}

	if len(b) > left {
		n, err = fb.readerAt.ReadAt(b[:left], fb.readed+fb.readRange.Begin)
	} else {
		n, err = fb.readerAt.ReadAt(b, fb.readed+fb.readRange.Begin)
	}

	n64 := int64(n)
	fb.readed += n64
	if fb.rateLimit != nil {
		fb.rateLimit.Add(n64) // 限速阻塞
	}
	if fb.speedsStatRef != nil {
		fb.speedsStatRef.Add(n64)
	}
	return
}

func (fb *fileBlock) Seek(offset int64, whence int) (int64, error) {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	switch whence {
	case os.SEEK_SET:
		fb.readed = offset
	case os.SEEK_CUR:
		fb.readed += offset
	case os.SEEK_END:
		fb.readed = fb.readRange.End - fb.readRange.Begin + offset
	default:
		return 0, fmt.Errorf("unsupport whence: %d", whence)
	}
	if fb.readed < 0 {
		fb.readed = 0
	}
	return fb.readed, nil
}

func (fb *fileBlock) Len() int64 {
	return fb.readRange.End - fb.readRange.Begin
}

func (fb *fileBlock) Left() int64 {
	return fb.readRange.End - fb.readRange.Begin - fb.readed
}

func (fb *fileBlock) Range() transfer.Range {
	return fb.readRange
}

func (fb *fileBlock) Readed() int64 {
	return fb.readed
}
