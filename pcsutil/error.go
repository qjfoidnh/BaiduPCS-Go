package pcsutil

import (
	"log"
	"os"
)

// PrintErrIfExist 简易错误处理, 如果 err 存在, 就只向屏幕输出 err 。
func PrintErrIfExist(err error) {
	if err != nil {
		log.Println(err)
	}
}

// PrintErrAndExit 简易错误处理, 如果 err 存在, 向屏幕输出 err 并退出, annotate 是加在 err 之前的注释信息。
func PrintErrAndExit(annotate string, err error) {
	if err != nil {
		log.Println(annotate, err)
		os.Exit(1)
	}
}
