package converter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// B byte
	B = (int64)(1 << (10 * iota))
	// KB kilobyte
	KB
	// MB megabyte
	MB
	// GB gigabyte
	GB
	// TB terabyte
	TB
	// PB petabyte
	PB
)

// ConvertFileSize 文件大小格式化输出
func ConvertFileSize(size int64, precision ...int) string {
	pint := "6"
	if len(precision) == 1 {
		pint = fmt.Sprint(precision[0])
	}
	if size < 0 {
		return "0B"
	}
	if size < KB {
		return fmt.Sprintf("%dB", size)
	}
	if size < MB {
		return fmt.Sprintf("%."+pint+"fKB", float64(size)/float64(KB))
	}
	if size < GB {
		return fmt.Sprintf("%."+pint+"fMB", float64(size)/float64(MB))
	}
	if size < TB {
		return fmt.Sprintf("%."+pint+"fGB", float64(size)/float64(GB))
	}
	if size < PB {
		return fmt.Sprintf("%."+pint+"fTB", float64(size)/float64(TB))
	}
	return fmt.Sprintf("%."+pint+"fPB", float64(size)/float64(PB))
}

// ParseFileSizeStr 将文件大小字符串转换成字节数
func ParseFileSizeStr(ss string) (size int64, err error) {
	if ss == "" {
		err = errors.New("converter: size is empty")
		return
	}
	if !(ss[0] == '.' || '0' <= ss[0] && ss[0] <= '9') {
		err = errors.New("converter: invalid size: " + ss)
		return
	}

	var i int
	for i = range ss[1:] {
		i++
		if ss[i] == '.' || ('0' <= ss[i] && ss[i] <= '9') {
			// 属于数字
			continue
		}
		break
	}
	if ss[i] == '.' || ('0' <= ss[i] && ss[i] <= '9') { // 最后一个分隔符是否为数字
		i++
	}

	var (
		sizeStr      = ss[:i] // 数字部分
		unitStr      = ss[i:] // 单位部分
		sizeFloat, _ = strconv.ParseFloat(sizeStr, 10)
	)
	switch strings.ToUpper(unitStr) {
	case "", "B":
		size = int64(sizeFloat)
	case "K", "KB":
		size = int64(sizeFloat * float64(KB))
	case "M", "MB":
		size = int64(sizeFloat * float64(MB))
	case "G", "GB":
		size = int64(sizeFloat * float64(GB))
	case "T", "TB":
		size = int64(sizeFloat * float64(TB))
	case "P", "PB":
		size = int64(sizeFloat * float64(PB))
	default:
		err = errors.New("converter: invalid unit " + unitStr)
	}
	return
}
