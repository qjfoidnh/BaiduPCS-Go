package checksum

import (
	"hash"
	"io"
)

type (
	ChecksumWriter interface {
		io.Writer
		Sum() interface{}
	}

	ChecksumWriteUnit struct {
		SliceEnd       int64
		End            int64
		SliceSum       interface{}
		Sum            interface{}
		OnlySliceSum   bool
		ChecksumWriter ChecksumWriter

		ptr int64
	}

	hashChecksumWriter struct {
		h hash.Hash
	}

	hash32ChecksumWriter struct {
		h hash.Hash32
	}
)

func (wi *ChecksumWriteUnit) handleEnd() error {
	if wi.ptr >= wi.End {
		// 已写完
		if !wi.OnlySliceSum {
			wi.Sum = wi.ChecksumWriter.Sum()
		}
		return ErrChecksumWriteStop
	}
	return nil
}

func (wi *ChecksumWriteUnit) write(p []byte) (n int, err error) {
	if wi.End <= 0 {
		// do nothing
		err = ErrChecksumWriteStop
		return
	}
	err = wi.handleEnd()
	if err != nil {
		return
	}

	var (
		i    int
		left = wi.End - wi.ptr
		lenP = len(p)
	)
	if left < int64(lenP) {
		// 读取即将完毕
		i = int(left)
	} else {
		i = lenP
	}
	n, err = wi.ChecksumWriter.Write(p[:i])
	if err != nil {
		return
	}
	wi.ptr += int64(n)
	if left < int64(lenP) {
		err = wi.handleEnd()
		return
	}
	return
}

func (wi *ChecksumWriteUnit) Write(p []byte) (n int, err error) {
	if wi.SliceEnd <= 0 { // 忽略Slice
		// 读取全部
		n, err = wi.write(p)
		return
	}

	// 要计算Slice的情况
	// 调整slice
	if wi.SliceEnd > wi.End {
		wi.SliceEnd = wi.End
	}

	// 计算剩余Slice
	var (
		sliceLeft = wi.SliceEnd - wi.ptr
	)
	if sliceLeft <= 0 {
		// 已处理完Slice
		if wi.OnlySliceSum {
			err = ErrChecksumWriteStop
			return
		}

		// 继续处理
		n, err = wi.write(p)
		return
	}

	var (
		lenP = len(p)
	)
	if sliceLeft <= int64(lenP) {
		var n1, n2 int
		n1, err = wi.write(p[:sliceLeft])
		n += n1
		if err != nil {
			return
		}
		wi.SliceSum = wi.ChecksumWriter.Sum().([]byte)
		n2, err = wi.write(p[sliceLeft:])
		n += n2
		if err != nil {
			return
		}
		return
	}
	n, err = wi.write(p)
	return
}

func NewHashChecksumWriter(h hash.Hash) ChecksumWriter {
	return &hashChecksumWriter{
		h: h,
	}
}

func (hc *hashChecksumWriter) Write(p []byte) (n int, err error) {
	return hc.h.Write(p)
}

func (hc *hashChecksumWriter) Sum() interface{} {
	return hc.h.Sum(nil)
}

func NewHash32ChecksumWriter(h32 hash.Hash32) ChecksumWriter {
	return &hash32ChecksumWriter{
		h: h32,
	}
}

func (hc *hash32ChecksumWriter) Write(p []byte) (n int, err error) {
	return hc.h.Write(p)
}

func (hc *hash32ChecksumWriter) Sum() interface{} {
	return hc.h.Sum32()
}
