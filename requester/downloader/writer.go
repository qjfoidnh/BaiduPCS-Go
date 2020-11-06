package downloader

import (
	"io"
	"os"
)

type (
	// Fder 获取fd接口
	Fder interface {
		Fd() uintptr
	}

	// Writer 下载器数据输出接口
	Writer interface {
		io.WriterAt
	}
)

// NewDownloaderWriterByFilename 创建下载器数据输出接口, 类似于os.OpenFile
func NewDownloaderWriterByFilename(name string, flag int, perm os.FileMode) (writer Writer, file *os.File, err error) {
	file, err = os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}

	writer = file
	return
}
