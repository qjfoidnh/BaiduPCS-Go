package pcscommand

import (
	"fmt"
        "strconv"	
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
		ShowSize bool
		ShowFsid bool
	}
)

func Format(n int64) string {
    in := strconv.FormatInt(n, 10)
    numOfDigits := len(in)
    if n < 0 {
        numOfDigits-- // First character is the - sign (not a digit)
    }
    numOfCommas := (numOfDigits - 1) / 3

    out := make([]byte, len(in)+numOfCommas)
    if n < 0 {
        in, out[0] = in[1:], '-'
    }

    for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
        out[j] = in[i]
        if i == 0 {
            return string(out)
        }
        if k++; k == 3 {
            j, k = j-1, 0
            out[j] = ','
        }
    }
}


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
	    if option.ShowFsid && option.ShowSize {
	      fmt.Printf("%v%v %v: %v: %v\n", indentPrefixStr, prefix, file.Filename, Format(file.Size), file.FsID)
	    } else if option.ShowFsid {
	      fmt.Printf("%v%v %v: %v\n", indentPrefixStr, prefix, file.Filename, file.FsID)
	    } else if option.ShowSize {
	      fmt.Printf("%v%v %v: %v\n", indentPrefixStr, prefix, file.Filename, Format(file.Size))
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
