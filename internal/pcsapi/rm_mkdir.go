package pcsapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

type RemoveStructure struct {
	TargetPaths []string `json:"target_paths" form:"target_paths" binding:"required"`
}

func remove(paths ...string) (fileList []string, err error) {
	fileList, err = matchPathByShellPattern(paths...)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pcscommand.GetBaiduPCS().Remove(paths...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("操作成功, 文件/目录已删除, 可在网盘文件回收站找回: ")
	return
}

func runRemove(ctx *gin.Context) {
	var args GetMetaStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("ls command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	fileList, err := remove(args.TargetPaths...)
	if err != nil {
		fmt.Printf("rm command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error":             err.Error(),
			"failed_files":      fileList,
			"failed_file_count": len(fileList),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success_files":       fileList,
		"success_files_count": len(fileList),
	})
}

// 将RunRemove挂载到路由列表中
func initRunRemove(group *gin.RouterGroup) {
	group.GET("rm", runRemove)
	group.POST("rm", runRemove)
}

type MakeDirectoryStructure struct {
	TargetPath string `json:"target_path" form:"target_path" binding:"required"`
}

func runMKdir(ctx *gin.Context) {
	var args MakeDirectoryStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("mkdr command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	activeUser := pcscommand.GetActiveUser()
	err := pcscommand.GetBaiduPCS().Mkdir(activeUser.PathJoin(args.TargetPath))
	if err != nil {
		fmt.Printf("mkdr command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "创建目录成功",
	})
}

func initRunMKdir(group *gin.RouterGroup) {
	group.GET("mkdir", runMKdir)
	group.POST("mkdir", runMKdir)
}
