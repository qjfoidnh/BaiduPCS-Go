package checksum

import (
	"bytes"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"path/filepath"
)

// EqualLengthMD5 检测md5和大小是否相同
func (lfm *LocalFileMeta) EqualLengthMD5(m *LocalFileMeta) bool {
	if lfm.Length != m.Length {
		return false
	}
	if bytes.Compare(lfm.MD5, m.MD5) != 0 {
		return false
	}
	return true
}

// CompleteAbsPath 补齐绝对路径
func (lfm *LocalFileMeta) CompleteAbsPath() {
	if filepath.IsAbs(lfm.Path) {
		return
	}

	absPath, err := filepath.Abs(lfm.Path)
	if err != nil {
		return
	}

	lfm.Path = absPath
}

// GetFileSum 获取文件的大小, md5, 前256KB切片的 md5, crc32
func GetFileSum(localPath string, flag int) (lfc *LocalFileChecksum, err error) {
	lfc = NewLocalFileChecksum(localPath, int(baidupcs.SliceMD5Size))
	defer lfc.Close()

	err = lfc.OpenPath()
	if err != nil {
		return nil, err
	}

	err = lfc.Sum(flag)
	if err != nil {
		return nil, err
	}
	return lfc, nil
}
