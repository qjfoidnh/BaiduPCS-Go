// Package converter 格式, 类型转换包
package converter

import (
	"github.com/mattn/go-runewidth"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

const (
	// InvalidChars 文件名中的非法字符
	InvalidChars = `\/:*?"<>|`
)

// ToString unsafe 转换, 将 []byte 转换为 string
func ToString(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

// ToBytes unsafe 转换, 将 string 转换为 []byte
func ToBytes(str string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}))
}

// ToBytesUnsafe unsafe 转换, 请确保转换后的 []byte 不涉及 cap() 操作, 将 string 转换为 []byte
func ToBytesUnsafe(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

// IntToBool int 类型转换为 bool
func IntToBool(i int) bool {
	return i != 0
}

// SliceInt64ToString []int64 转换为 []string
func SliceInt64ToString(si []int64) (ss []string) {
	ss = make([]string, 0, len(si))
	for k := range si {
		ss = append(ss, strconv.FormatInt(si[k], 10))
	}
	return ss
}

// SliceStringToInt64 []string 转换为 []int64
func SliceStringToInt64(ss []string) (si []int64) {
	si = make([]int64, 0, len(ss))
	var (
		i   int64
		err error
	)
	for k := range ss {
		i, err = strconv.ParseInt(ss[k], 10, 64)
		if err != nil {
			continue
		}
		si = append(si, i)
	}
	return
}

// SliceStringToInt []string 转换为 []int
func SliceStringToInt(ss []string) (si []int) {
	si = make([]int, 0, len(ss))
	var (
		i   int
		err error
	)
	for k := range ss {
		i, err = strconv.Atoi(ss[k])
		if err != nil {
			continue
		}
		si = append(si, i)
	}
	return
}

// MustInt 将string转换为int, 忽略错误
func MustInt(s string) (n int) {
	n, _ = strconv.Atoi(s)
	return
}

// MustInt64 将string转换为int64, 忽略错误
func MustInt64(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return
}

// ShortDisplay 缩略显示字符串s, 显示长度为num, 缩略的内容用"..."填充
func ShortDisplay(s string, num int) string {
	var (
		sb = strings.Builder{}
		n  int
	)
	for _, v := range s {
		if unicode.Is(unicode.C, v) { // 去除无效字符
			continue
		}
		n += runewidth.RuneWidth(v)
		if n > num {
			sb.WriteString("...")
			break
		}
		sb.WriteRune(v)
	}

	return sb.String()
}

// TrimPathInvalidChars 清除文件名中的非法字符
func TrimPathInvalidChars(fpath string) string {
	buf := make([]byte, 0, len(fpath))

	for _, c := range ToBytesUnsafe(fpath) {
		if strings.ContainsRune(InvalidChars, rune(c)) {
			continue
		}

		buf = append(buf, c)
	}

	return ToString(buf)
}
