package pcsapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
)

type ListDirectoryStructure struct {
	Asc     bool   `form:"asc,omitempty" json:"asc,omitempty"`
	Desc    bool   `form:"desc,omitempty" json:"desc,omitempty"`
	Time    bool   `form:"time,omitempty" json:"time,omitempty"`
	Name    bool   `form:"name,omitempty" json:"name,omitempty"`
	Size    bool   `form:"size,omitempty" json:"size,omitempty"`
	Pcspath string `form:"pcspath" json:"pcspath"`
}

func runListDirectory(ctx *gin.Context) {
	var args = ListDirectoryStructure{
		Asc:     true,
		Name:    true,
		Pcspath: ".",
	}
	orderOptions := &baidupcs.OrderOptions{}
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("ls command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	switch {
	case args.Asc:
		orderOptions.Order = baidupcs.OrderAsc
	case args.Desc:
		orderOptions.Order = baidupcs.OrderDesc
	default:
		orderOptions.Order = baidupcs.OrderAsc
	}
	switch {
	case args.Time:
		orderOptions.By = baidupcs.OrderByTime
	case args.Name:
		orderOptions.By = baidupcs.OrderByName
	case args.Size:
		orderOptions.By = baidupcs.OrderBySize
	default:
		orderOptions.By = baidupcs.OrderByName
	}
	files, err := runLs(args.Pcspath, &pcscommand.LsOptions{
		Total: true,
	}, orderOptions)
	if err != nil {
		fmt.Printf("ls command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"files":     files,
		"fileCount": len(files),
	})
}

func initRunListDirectory(group *gin.RouterGroup) {
	// get方式下的pcspath参数中的+会和url编码后的空格混淆
	// group.GET("ls", runListDirectory)
	group.POST("ls", runListDirectory)
}

type SearchFilesStructure struct {
	Path    string `json:"path" form:"path"`
	Recurse bool   `form:"recurse" json:"recurse"`
	Keyword string `json:"keyword" form:"keyword" binding:"required"`
}

func runSearchFiles(ctx *gin.Context) {
	var args SearchFilesStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("search command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := matchPathByShellPatternOnce(&args.Path)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	files, err := pcscommand.GetBaiduPCS().Search(args.Path, args.Keyword, args.Recurse)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"files":     files,
		"fileCount": len(files),
	})
}

// 将RunSearchFiles挂载到路由列表中
func initRunSearchFiles(group *gin.RouterGroup) {
	group.GET("search", runSearchFiles)
	group.POST("search", runSearchFiles)
}

func runPWD(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"workingDir": pcsconfig.Config.ActiveUser().Workdir,
	})
}

// 将RunPWD挂载到路由列表中
func initRunPWD(group *gin.RouterGroup) {
	group.GET("pwd", runPWD)
	group.POST("pwd", runPWD)
}
