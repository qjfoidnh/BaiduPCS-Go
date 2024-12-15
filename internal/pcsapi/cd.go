package pcsapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
)

type ChangeDirectoryStructure struct {
	TargetPath string `form:"target_path" json:"target_path" binding:"required"`
	IsList     bool   `form:"is_list,omitempty" json:"is_list,omitempty"`
}

func changeDirectory(targetPath string, isList bool) (files baidupcs.FileDirectoryList, err error) {
	pcs := pcscommand.GetBaiduPCS()
	err = matchPathByShellPatternOnce(&targetPath)
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
		err = fmt.Errorf("错误: %s 不是一个目录 (文件夹)", targetPath)
		fmt.Println(err.Error())
		return
	}

	pcscommand.GetActiveUser().Workdir = targetPath
	pcsconfig.Config.Save()

	fmt.Printf("改变工作目录: %s\n", targetPath)

	if isList {
		files, err = pcs.FilesDirectoriesList(targetPath, nil)
	}
	return
}

func runChangeDirectory(ctx *gin.Context) {
	var args ChangeDirectoryStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("cd command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	files, err := changeDirectory(args.TargetPath, args.IsList)
	rep := gin.H{}
	if err != nil {
		rep["error"] = err
	}
	if files != nil {
		rep["current_files"] = files
	}
	ctx.JSON(http.StatusAccepted, rep)
}

// 将RunChangeDirectory函数挂载到路由列表中
func initRunChangeDirectory(group *gin.RouterGroup) {
	group.POST("cd", runChangeDirectory)
}
