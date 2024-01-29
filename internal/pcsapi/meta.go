package pcsapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

type GetMetaStructure struct {
	TargetPaths []string `json:"target_paths" form:"target_paths"`
}

func GetMeta(targetPaths []string) (fileList []map[string]any, err error) {
	if len(targetPaths) == 0 {
		targetPaths = append(targetPaths, ".")
	}
	targetPaths, err = matchPathByShellPattern(targetPaths...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileList_ := make(baidupcs.FileDirectoryList, len(targetPaths))
	for i, v := range targetPaths {
		var data *baidupcs.FileDirectory
		data, err = pcscommand.GetBaiduPCS().FilesDirectoriesMeta(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileList_[i] = data
	}
	fileList = renderFileList(fileList_)
	return
}

func runGetMeta(ctx *gin.Context) {
	var args GetMetaStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("ls command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	fileList, err := GetMeta(args.TargetPaths)
	if err != nil {
		fmt.Printf("meta command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"metas": fileList,
	})
}

// 将RunGetMeta挂载到路由列表中
func initRunGetMeta(group *gin.RouterGroup) {
	// group.GET("meta", runGetMeta)
	group.POST("meta", runGetMeta)
}
