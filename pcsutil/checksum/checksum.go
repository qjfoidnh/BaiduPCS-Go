// Package checksum 校验本地文件包
package checksum

import (
	"crypto/md5"
	"github.com/iikira/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"hash/crc32"
	"io"
	"os"
)

const (
	// DefaultBufSize 默认的bufSize
	DefaultBufSize = int(256 * converter.KB)
)

const (
	// CHECKSUM_MD5 获取文件的 md5 值
	CHECKSUM_MD5 int = 1 << iota
	// CHECKSUM_SLICE_MD5 获取文件前 sliceSize 切片的 md5 值
	CHECKSUM_SLICE_MD5
	// CHECKSUM_CRC32 获取文件的 crc32 值
	CHECKSUM_CRC32
)

type (
	// LocalFileMeta 本地文件元信息
	LocalFileMeta struct {
		Path     string `json:"path"`     // 本地路径
		Length   int64  `json:"length"`   // 文件大小
		SliceMD5 []byte `json:"slicemd5"` // 文件前 requiredSliceLen (256KB) 切片的 md5 值
		MD5      []byte `json:"md5"`      // 文件的 md5
		CRC32    uint32 `json:"crc32"`    // 文件的 crc32
		ModTime  int64  `json:"modtime"`  // 修改日期
	}

	// LocalFileChecksum 校验本地文件
	LocalFileChecksum struct {
		LocalFileMeta
		bufSize   int
		sliceSize int
		buf       []byte
		file      *os.File // 文件
	}
)

func NewLocalFileChecksum(localPath string, sliceSize int) *LocalFileChecksum {
	return NewLocalFileChecksumWithBufSize(localPath, DefaultBufSize, sliceSize)
}

func NewLocalFileChecksumWithBufSize(localPath string, bufSize, sliceSize int) *LocalFileChecksum {
	return &LocalFileChecksum{
		LocalFileMeta: LocalFileMeta{
			Path: localPath,
		},
		bufSize:   bufSize,
		sliceSize: sliceSize,
	}
}

// OpenPath 检查文件状态并获取文件的大小 (Length)
func (lfc *LocalFileChecksum) OpenPath() error {
	if lfc.file != nil {
		lfc.file.Close()
	}

	var err error
	lfc.file, err = os.Open(lfc.Path)
	if err != nil {
		return err
	}

	info, err := lfc.file.Stat()
	if err != nil {
		return err
	}

	lfc.Length = info.Size()
	lfc.ModTime = info.ModTime().Unix()
	return nil
}

// GetFile 获取文件
func (lfc *LocalFileChecksum) GetFile() *os.File {
	return lfc.file
}

// Close 关闭文件
func (lfc *LocalFileChecksum) Close() error {
	if lfc.file == nil {
		return ErrFileIsNil
	}

	return lfc.file.Close()
}

func (lfc *LocalFileChecksum) initBuf() {
	if lfc.buf == nil {
		lfc.buf = cachepool.RawMallocByteSlice(lfc.bufSize)
	}
}

func (lfc *LocalFileChecksum) writeChecksum(data []byte, wus ...*ChecksumWriteUnit) (err error) {
	doneCount := 0
	for _, wu := range wus {
		_, err := wu.Write(data)
		switch err {
		case ErrChecksumWriteStop:
			doneCount++
			continue
		case nil:
		default:
			return err
		}
	}
	if doneCount == len(wus) {
		return ErrChecksumWriteAllStop
	}
	return nil
}

func (lfc *LocalFileChecksum) repeatRead(wus ...*ChecksumWriteUnit) (err error) {
	if lfc.file == nil {
		return ErrFileIsNil
	}

	lfc.initBuf()

	defer func() {
		_, err = lfc.file.Seek(0, os.SEEK_SET) // 恢复文件指针
		if err != nil {
			return
		}
	}()

	// 读文件
	var (
		n int
	)
read:
	for {
		n, err = lfc.file.Read(lfc.buf)
		switch err {
		case io.EOF:
			err = lfc.writeChecksum(lfc.buf[:n], wus...)
			break read
		case nil:
			err = lfc.writeChecksum(lfc.buf[:n], wus...)
		default:
			return
		}
	}
	switch err {
	case ErrChecksumWriteAllStop: // 全部结束
		err = nil
	}
	return
}

func (lfc *LocalFileChecksum) createChecksumWriteUnit(cw ChecksumWriter, isAll, isSlice bool, getSumFunc func(sliceSum interface{}, sum interface{})) (wu *ChecksumWriteUnit, deferFunc func(err error)) {
	wu = &ChecksumWriteUnit{
		ChecksumWriter: cw,
		End:            lfc.LocalFileMeta.Length,
		OnlySliceSum:   !isAll,
	}

	if isSlice {
		wu.SliceEnd = int64(lfc.sliceSize)
	}

	return wu, func(err error) {
		if err != nil {
			return
		}
		getSumFunc(wu.SliceSum, wu.Sum)
	}
}

// Sum 计算文件摘要值
func (lfc *LocalFileChecksum) Sum(checkSumFlag int) (err error) {
	lfc.fix()
	wus := make([]*ChecksumWriteUnit, 0, 2)
	if (checkSumFlag & (CHECKSUM_MD5 | CHECKSUM_SLICE_MD5)) != 0 {
		md5w := md5.New()
		wu, d := lfc.createChecksumWriteUnit(
			NewHashChecksumWriter(md5w),
			(checkSumFlag&CHECKSUM_MD5) != 0,
			(checkSumFlag&CHECKSUM_SLICE_MD5) != 0,
			func(sliceSum interface{}, sum interface{}) {
				if sliceSum != nil {
					lfc.SliceMD5 = sliceSum.([]byte)
				}
				if sum != nil {
					lfc.MD5 = sum.([]byte)
				}
			},
		)

		wus = append(wus, wu)
		defer d(err)
	}
	if (checkSumFlag & CHECKSUM_CRC32) != 0 {
		crc32w := crc32.NewIEEE()
		wu, d := lfc.createChecksumWriteUnit(
			NewHash32ChecksumWriter(crc32w),
			true,
			false,
			func(sliceSum interface{}, sum interface{}) {
				if sum != nil {
					lfc.CRC32 = sum.(uint32)
				}
			},
		)

		wus = append(wus, wu)
		defer d(err)
	}

	err = lfc.repeatRead(wus...)
	return
}

func (lfc *LocalFileChecksum) fix() {
	if lfc.sliceSize <= 0 {
		lfc.sliceSize = DefaultBufSize
	}
	if lfc.bufSize < DefaultBufSize {
		lfc.bufSize = DefaultBufSize
	}
}
