package pcscommand

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"strings"
)

const (
	indentPrefix   = "│   "
	pathPrefix     = "├──"
	lastFilePrefix = "└──"
)

type (
	TreeOptions struct {
		Depth    int
		ShowFsid bool
	}
)

func getTree(pcspath string, depth int, option *TreeOptions) {
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
			if option.ShowFsid {
				fmt.Printf("%v%v %v/: %v\n", indentPrefixStr, pathPrefix, file.Filename, file.FsID)
			} else {
				fmt.Printf("%v%v %v/\n", indentPrefixStr, pathPrefix, file.Filename)
			}
			if option.Depth < 0 || depth < option.Depth {
				getTree(file.Path, depth+1, option)
			}
			continue
		}

		if i+1 == fN {
			prefix = lastFilePrefix
		}
		if option.ShowFsid {
			fmt.Printf("%v%v %v: %v\n", indentPrefixStr, prefix, file.Filename, file.FsID)
		} else {
			fmt.Printf("%v%v %v\n", indentPrefixStr, prefix, file.Filename)
		}
	}

	return
}

// RunTree 列出树形图
func RunTree(path string, depth int, option *TreeOptions) {
	getTree(path, depth, option)
}
