package pcscommand

import (
	"encoding/json"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

type MetaListJsonObject struct {
	err  *string
	data *baidupcs.FileDirectoryList
}

// RunGetMeta 执行 获取文件/目录的元信息
func RunGetMeta(jsonOutput bool, targetPaths ...string) {
	targetPaths, err := matchPathByShellPattern(targetPaths...)
	if printErrorAndReturn(jsonOutput, err) {
		return
	}

	for k, targetPath := range targetPaths {
		if !jsonOutput {
			fmt.Printf("[%d] - [%s] --------------\n", k, targetPath)
		}
		data, err := GetBaiduPCS().FilesDirectoriesMeta(targetPath)
		if printErrorAndReturn(jsonOutput, err) {
			return
		}
		if !jsonOutput {
			fmt.Println()
			fmt.Println(data)
		} else {
			byteStr, err := json.Marshal(data)
			if printErrorAndReturn(jsonOutput, err) {
				return
			}
			println(string(byteStr))
		}
	}
}

func printErrorAndReturn(jsonOutput bool, err error) bool {
	result := MetaListJsonObject{}
	if err != nil {
		errorString := err.Error()
		if jsonOutput {
			result.err = &errorString
			jsonByte, _ := json.Marshal(result)
			errorString = string(jsonByte)
		}
		fmt.Println(errorString)
		return true
	}
	return false
}
