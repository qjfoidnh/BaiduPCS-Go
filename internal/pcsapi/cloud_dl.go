package pcsapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

type CloudDlAddTaskStructure struct {
	SourceURLs []string `form:"sourceURLs" json:"sourceURLs" binding:"required"`
	SavePath   string   `form:"savePath" json:"savePath" binding:"required"`
}

func runCloudDlAddTask(ctx *gin.Context) {
	var args CloudDlAddTaskStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("offlinedl command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	err, successids, failedids := CloudDlAddTask(args.SourceURLs, args.SavePath)
	rep := gin.H{
		"successIds": successids,
		"failedIds":  failedids,
	}
	if err != nil {
		rep["error"] = err
	}
	ctx.JSON(http.StatusOK, rep)
}

func CloudDlAddTask(sourceURLs []string, savePath string) (err error, successIds []int64, failedIds []int64) {
	pcs := pcscommand.GetBaiduPCS()
	err = matchPathByShellPatternOnce(&savePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var taskid int64
	for k := range sourceURLs {
		taskid, err = pcs.CloudDlAddTask(sourceURLs[k], savePath+baidupcs.PathSeparator)
		if err != nil {
			fmt.Printf("[%d] %s, 地址: %s\n", k+1, err, sourceURLs[k])
			failedIds = append(failedIds, taskid)
			continue
		}
		successIds = append(successIds, taskid)
		fmt.Printf("[%d] 添加离线任务成功, 任务ID(task_id): %d, 源地址: %s, 保存路径: %s\n", k+1, taskid, sourceURLs[k], savePath)
	}
	return
}

// 挂载RunCloudDlAddTask函数
func initRunCloudDlAddTask(group *gin.RouterGroup) {
	group.POST("offlinedl", runCloudDlAddTask)
}

type CloudDlQueryTaskStructure struct {
	TaskIDs []int64 `form:"taskIDs" json:"taskIDs" binding:"required"`
}

func runCloudDlQueryTask(ctx *gin.Context) {
	var args CloudDlQueryTaskStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("query command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	cl, err := pcscommand.GetBaiduPCS().CloudDlQueryTask(args.TaskIDs)
	rep := gin.H{}
	rep["taskInfos"] = cl
	if err != nil {
		rep["error"] = err
	}
	ctx.JSON(http.StatusOK, rep)
}

// 挂载RunCloudDlQueryTask函数
func initRunCloudDlQueryTask(group *gin.RouterGroup) {
	group.GET("query", runCloudDlQueryTask)
	group.POST("query", runCloudDlQueryTask)
}

func runCloudDlListTask(ctx *gin.Context) {
	cl, err := pcscommand.GetBaiduPCS().CloudDlListTask()
	rep := gin.H{}
	rep["taskInfos"] = cl
	if err != nil {
		rep["error"] = err
	}
	ctx.JSON(http.StatusOK, cl)
}

// 挂载RunCloudDlListTask函数
func initRunCloudDlListTask(group *gin.RouterGroup) {
	group.GET("list", runCloudDlListTask)
	group.POST("list", runCloudDlListTask)
}

type CloudDlDeleteTaskStructure struct {
	TaskIDs []int64 `form:"taskIDs" json:"taskIDs" binding:"required"`
}

func runCloudDlDeleteTask(ctx *gin.Context) {
	var args CloudDlDeleteTaskStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("delete command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	var errs []error
	var successIds []int64
	for _, id := range args.TaskIDs {
		err := pcscommand.GetBaiduPCS().CloudDlDeleteTask(id)
		if err != nil {
			fmt.Printf("[%d] %s\n", id, err)
			errs = append(errs, err)
			continue
		}
		successIds = append(successIds, id)
		fmt.Printf("[%d] 删除成功\n", id)
	}
	if len(errs) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"errors": errs,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"successIds": successIds,
	})
}

// 挂载RunCloudDlDeleteTask函数
func initRunCloudDlDeleteTask(group *gin.RouterGroup) {
	group.GET("delete", runCloudDlDeleteTask)
	group.POST("delete", runCloudDlDeleteTask)
}

type CloudDlCancelTaskStructure CloudDlDeleteTaskStructure

func runCloudDlCancelTask(ctx *gin.Context) {
	var args CloudDlCancelTaskStructure
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("cancel command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(args.TaskIDs) <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"error": fmt.Errorf("no task IDs"),
		})
		return
	}
	var (
		errs       []error
		successIds []int64
	)
	for _, id := range args.TaskIDs {
		err := pcscommand.GetBaiduPCS().CloudDlCancelTask(id)
		if err != nil {
			errs = append(errs, err)
		} else {
			successIds = append(successIds, id)
		}
	}
	rep := gin.H{}
	if len(errs) > 0 {
		rep["error"] = errs
	}
	if len(successIds) > 0 {
		rep["successIds"] = successIds
	}
	ctx.JSON(http.StatusOK, rep)
}

// 挂载RunCloudDlClearTask函数
func initRunCloudDlCancelTask(group *gin.RouterGroup) {
	group.GET("cancel", runCloudDlCancelTask)
	group.POST("cancel", runCloudDlCancelTask)
}

func initRunOffLineDownload(router *gin.RouterGroup) {
	offlineDL := router.Group("offelinedl")
	initRunCloudDlAddTask(offlineDL)
	initRunCloudDlQueryTask(offlineDL)
	initRunCloudDlListTask(offlineDL)
	initRunCloudDlCancelTask(offlineDL)
	initRunCloudDlDeleteTask(offlineDL)
}
