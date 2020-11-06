package pcsutil

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/iikira/BaiduPCS-Go/pcsutil/pcstime"
	"log"
)

var (
	// ErrorColor 设置输出错误的颜色
	ErrorColor = color.New(color.FgRed).SprintFunc()
)

// 自定义log writer
type logWriter struct{}

func (logWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprint(color.Output, "["+pcstime.BeijingTimeOption("Refer")+"] "+string(bytes))
}

// SetLogPrefix 设置日志输出的时间前缀
func SetLogPrefix() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}
