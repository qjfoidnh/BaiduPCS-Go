package pcscommand

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

// RunFixMD5 执行修复md5
func RunFixMD5(pcspaths ...string) {
	absPaths, err := matchPathByShellPattern(pcspaths...)
	if err != nil {
		fmt.Println(err)
		return
	}

	pcs := GetBaiduPCS()
	finfoList, err := pcs.FilesDirectoriesBatchMeta(absPaths...)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, finfo := range finfoList {
		err := pcs.FixMD5ByFileInfo(finfo)
		if err == nil {
			fmt.Printf("[%d] - [%s] 修复md5成功\n", k, finfo.Path)
			continue
		}

		if err.GetError() == baidupcs.ErrFixMD5Failed {
			fmt.Printf("[%d] - [%s] 修复md5失败, 可能是服务器未刷新\n", k, finfo.Path)
			continue
		}
		fmt.Printf("[%d] - [%s] 修复md5失败, 错误信息: %s\n", k, finfo.Path, err)
	}
}
