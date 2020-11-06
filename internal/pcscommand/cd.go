package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
)

// RunChangeDirectory 执行更改工作目录
func RunChangeDirectory(targetPath string, isList bool) {
	pcs := GetBaiduPCS()
	err := matchPathByShellPatternOnce(&targetPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := pcs.FilesDirectoriesMeta(targetPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !data.Isdir {
		fmt.Printf("错误: %s 不是一个目录 (文件夹)\n", targetPath)
		return
	}

	GetActiveUser().Workdir = targetPath
	pcsconfig.Config.Save()

	fmt.Printf("改变工作目录: %s\n", targetPath)

	if isList {
		RunLs(".", nil, baidupcs.DefaultOrderOptions)
	}
}
