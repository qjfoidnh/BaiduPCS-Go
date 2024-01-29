package pcsapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

func runLs(pcspath string, lsOptions *pcscommand.LsOptions, orderOptions *baidupcs.OrderOptions) (files []map[string]any, err error) {
	err = matchPathByShellPatternOnce(&pcspath)
	if err != nil {
		fmt.Println(err)
		return
	}
	var files_ baidupcs.FileDirectoryList
	files_, err = pcscommand.GetBaiduPCS().FilesDirectoriesList(pcspath, orderOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("\n当前目录: %s\n----\n", pcspath)

	if lsOptions == nil {
		lsOptions = &pcscommand.LsOptions{}
	}

	files = renderFileList(files_)
	return
}

func renderFileList(files baidupcs.FileDirectoryList) (fileList []map[string]any) {
	fileList = make([]map[string]any, len(files))
	for i, v := range files {
		fileList[i] = gin.H{
			"FsID":        v.FsID,
			"AppID":       v.AppID,
			"Path":        v.Path,
			"Filename":    v.Filename,
			"Ctime":       v.Ctime,
			"Mtime":       v.Mtime,
			"MD5":         v.MD5,
			"Size":        v.Size,
			"Isdir":       v.Isdir,
			"Ifhassubdir": v.Ifhassubdir,
			"PreBase":     v.PreBase,
		}
	}
	return
}

func RunSearch(targetPath string, keyward string, opt *pcscommand.SearchOptions) (files []map[string]any, err error) {
	err = matchPathByShellPatternOnce(&targetPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if opt == nil {
		opt = &pcscommand.SearchOptions{}
	}
	var files_ baidupcs.FileDirectoryList
	files_, err = pcscommand.GetBaiduPCS().Search(targetPath, keyward, opt.Recurse)
	if err != nil {
		fmt.Println(err)
		return
	}
	files = renderFileList(files_)
	return
}
