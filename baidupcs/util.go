package baidupcs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"io"
	"path"
	"strings"
	"time"
)

// Isdir 检查路径在网盘中是否为目录
func (pcs *BaiduPCS) Isdir(pcspath string) (fileSize int64, isdir bool, pcsError pcserror.Error) {
	if path.Clean(pcspath) == PathSeparator {
		return 0, true, nil
	}

	f, pcsError := pcs.FilesDirectoriesMeta(pcspath)
	if pcsError != nil {
		return 0, false, pcsError
	}

	return f.Size, f.Isdir, nil
}

func (pcs *BaiduPCS) CheckIsdir(op string, targetPath string, policy string, fileSize int64) pcserror.Error {
	// 检测文件是否存在于网盘路径
	// 很重要, 如果文件存在会直接覆盖!!! 即使是根目录!
	onlineSize, isdir, pcsError := pcs.Isdir(targetPath)
	if pcsError != nil {
		// 忽略远程服务端返回的错误
		if pcsError.GetErrType() != pcserror.ErrTypeRemoteError {
			return pcsError
		}
	}

	errInfo := pcserror.NewPCSErrorInfo(op)
	if isdir {
		errInfo.ErrType = pcserror.ErrTypeOthers
		errInfo.Err = errors.New("保存路径不可以覆盖目录")
		return errInfo
	}
	// 如果存在文件, 则根据upload策略选择返回的错误码
	if pcsError == nil {
		switch policy {
		case "fail":
			errInfo.ErrCode = 114514
			errInfo.ErrType = pcserror.ErrTypeRemoteError
			errInfo.ErrMsg = "目标位置存在同名文件"
			return errInfo
		case "skip":
			errInfo.ErrCode = 114514
			errInfo.ErrMsg = "目标位置存在同名文件"
			errInfo.ErrType = pcserror.ErrTypeRemoteError
			return errInfo
		case "rsync":
			if onlineSize == fileSize {
				errInfo.ErrCode = 1919810
				errInfo.ErrMsg = "目标位置文件大小与源文件一致"
				errInfo.ErrType = pcserror.ErrTypeRemoteError
				return errInfo
			}
		default:
			return nil
		}
	}
	return nil
}

func mergeStringList(a ...string) string {
	s := strings.Join(a, `","`)
	return `["` + s + `"]`
}

func mergeInt64List(si ...int64) string {
	i := converter.SliceInt64ToString(si)
	s := strings.Join(i, ",")
	return "[" + s + "]"
}

func allRelatedDir(pcspaths []string) (dirs []string) {
	for _, pcspath := range pcspaths {
		pathDir := path.Dir(pcspath)
		if !pcsutil.ContainsString(dirs, pathDir) {
			dirs = append(dirs, pathDir)
		}
	}
	return
}

func CreatePasswd() string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "Asswecan")
	io.WriteString(h, t.String())
	passwd := fmt.Sprintf("%x", h.Sum(nil))
	return passwd[0:4]
}

// GetHTTPScheme 获取 http 协议, https 或 http
func GetHTTPScheme(https bool) (scheme string) {
	if https {
		return "https"
	}
	return "http"
}
