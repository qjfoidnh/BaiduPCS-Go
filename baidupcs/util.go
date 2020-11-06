package baidupcs

import (
	"errors"
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"path"
	"strings"
)

// Isdir 检查路径在网盘中是否为目录
func (pcs *BaiduPCS) Isdir(pcspath string) (isdir bool, pcsError pcserror.Error) {
	if path.Clean(pcspath) == PathSeparator {
		return true, nil
	}

	f, pcsError := pcs.FilesDirectoriesMeta(pcspath)
	if pcsError != nil {
		return false, pcsError
	}

	return f.Isdir, nil
}

func (pcs *BaiduPCS) checkIsdir(op string, targetPath string) pcserror.Error {
	// 检测文件是否存在于网盘路径
	// 很重要, 如果文件存在会直接覆盖!!! 即使是根目录!
	isdir, pcsError := pcs.Isdir(targetPath)
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

// GetHTTPScheme 获取 http 协议, https 或 http
func GetHTTPScheme(https bool) (scheme string) {
	if https {
		return "https"
	}
	return "http"
}
