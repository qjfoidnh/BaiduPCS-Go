package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"strings"
)

const (
	indentPrefix   = "│   "
	pathPrefix     = "├──"
	lastFilePrefix = "└──"
)

func getTree(pcspath string, depth int) {
	var (
		err   error
		files baidupcs.FileDirectoryList
	)
	if depth == 0 {
		err := matchPathByShellPatternOnce(&pcspath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	files, err = GetBaiduPCS().FilesDirectoriesList(pcspath, baidupcs.DefaultOrderOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		prefix          = pathPrefix
		fN              = len(files)
		indentPrefixStr = strings.Repeat(indentPrefix, depth)
	)
	for i, file := range files {
		if file.Isdir {
			fmt.Printf("%v%v %v/\n", indentPrefixStr, pathPrefix, file.Filename)
			getTree(file.Path, depth+1)
			continue
		}

		if i+1 == fN {
			prefix = lastFilePrefix
		}

		fmt.Printf("%v%v %v\n", indentPrefixStr, prefix, file.Filename)
	}

	return
}

// RunTree 列出树形图
func RunTree(path string) {
	getTree(path, 0)
}
