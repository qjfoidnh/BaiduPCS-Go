package baidupcs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"io"
	"math/rand"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	SkipPolicy      = "skip"
	OverWritePolicy = "overwrite"
	RsyncPolicy     = "rsync"
)

func (pcs *BaiduPCS) policyTortype(policy string) string {
	switch policy {
	case SkipPolicy:
		return "2"
	}

	// 兜底为覆盖逻辑
	return "3"
}

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
	targetFileSize, isdir, pcsError := pcs.Isdir(targetPath)
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
		case SkipPolicy:
			errInfo.ErrCode = 114514
			errInfo.ErrType = pcserror.ErrTypeRemoteError
			errInfo.ErrMsg = "目标位置存在同名文件"
			return errInfo
		case RsyncPolicy:
			if targetFileSize == fileSize {
				errInfo.ErrCode = 1919810
				errInfo.ErrType = pcserror.ErrTypeRemoteError
				errInfo.ErrMsg = "目标位置存在相同文件"
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

func sortBlockList(checksumMap map[int]string) []string {
	keys := make([]int, 0, len(checksumMap))
	for k := range checksumMap {
		keys = append(keys, k)
	}
	sort.Ints(keys) // 升序排序

	// 2. 按排序后的 Key 提取 Value
	result := make([]string, 0, len(checksumMap))
	for _, k := range keys {
		result = append(result, checksumMap[k])
	}
	return result
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

func DecryptMD5(rawMD5 string) string {
	if len(rawMD5) != 32 {
		return rawMD5
	}
	var keychar string = rawMD5[9:10]
	match, _ := regexp.MatchString("[a-f0-9]", keychar)
	if match {
		return rawMD5
	}
	sliceFirst := fmt.Sprintf("%x", []rune(rawMD5)[9]-'g')
	sliceSecond := rawMD5[0:9] + sliceFirst + rawMD5[10:]
	sliceThird := ""
	for i := 0; i < len(sliceSecond); i++ {
		if sliceSecond[i:i+1] == "-" {
			sliceThird += fmt.Sprintf("%x", 15&i)
			continue
		}
		num, err := strconv.ParseInt(sliceSecond[i:i+1], 16, 64)
		if err != nil {
			return rawMD5
		}
		sliceThird += fmt.Sprintf("%x", int(num)^(15&i))
	}
	return sliceThird[8:16] + sliceThird[0:8] + sliceThird[24:32] + sliceThird[16:24]
}

func RandomElement[T any](s []T) T {
	if len(s) == 0 {
		var zero T // 对于空slice，返回类型的零值
		return zero
	}
	return s[rand.Intn(len(s))]
}
