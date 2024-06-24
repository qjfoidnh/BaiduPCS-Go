package pcsapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

type RecycleListStructure struct {
	Page int  `json:"page,omitempty" form:"page,omitempty"`
	All  bool `json:"all,omitempty" form:"all,omitempty"`
}

func runRecycleList(ctx *gin.Context) {

	args := RecycleListStructure{
		Page: 1,
	}

	if err := ctx.ShouldBind(&args); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	pcs := pcscommand.GetBaiduPCS()
	var (
		recycle_file_list baidupcs.RecycleFDInfoList
		errs              []error
	)
	for {
		fdl, err := pcs.RecycleList(args.Page)
		if err != nil {
			errs = append(errs, err)
		}
		if len(fdl) <= 0 {
			break
		} else {
			recycle_file_list = append(recycle_file_list, fdl...)
		}
		if !args.All {
			break
		}
	}
	rep := gin.H{
		"recycle_file_list": recycle_file_list,
	}
	if len(errs) > 0 {
		rep["error"] = errs
	}
	ctx.JSON(http.StatusOK, rep)
}

func initRunRecycleList(group *gin.RouterGroup) {
	group.GET("list", runRecycleList)
	group.POST("list", runRecycleList)
}

type RecycleRestoreStructure struct {
	FidList []int64 `json:"fid_list,omitempty" form:"fid_list,omitempty" binding:"required"`
}

func runRecycleRestore(ctx *gin.Context) {
	args := RecycleRestoreStructure{}
	if err := ctx.ShouldBind(&args); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	var (
		pcs     = pcscommand.GetBaiduPCS()
		ex, err = pcs.RecycleRestore(args.FidList...)
	)
	rep := gin.H{
		"success_fids": ex,
	}
	if err != nil {
		rep["error"] = err.GetError()
	}
	ctx.JSON(http.StatusOK, rep)
}

func initRunRecycleRestore(router *gin.RouterGroup) {
	router.GET("restore", runRecycleRestore)
	router.POST("restore", runRecycleRestore)
}

func runRecycleClear(ctx *gin.Context) {
	pcs := pcscommand.GetBaiduPCS()
	susSum, err := pcs.RecycleClear()
	rep := gin.H{
		"success_nums": susSum,
	}
	if err != nil {
		rep["error"] = err.GetError()
	}
	ctx.JSON(http.StatusOK, rep)
}

func initRunRecycleClear(router *gin.RouterGroup) {
	router.POST("clear", runRecycleClear)
	router.GET("clear", runRecycleClear)
}

type RecycleDeleteStructure RecycleRestoreStructure

func runRecycleDelete(ctx *gin.Context) {
	args := RecycleDeleteStructure{}
	if err := ctx.ShouldBind(&args); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	pcs := pcscommand.GetBaiduPCS()
	err := pcs.RecycleDelete(args.FidList...)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.GetError(),
		})
	} else {
		ctx.JSON(http.StatusAccepted, nil)
	}
}

func initRunRecycleDelete(router *gin.RouterGroup) {
	router.POST("delete", runRecycleDelete)
	router.GET("delete", runRecycleDelete)
}

func initRunRecycle(router *gin.RouterGroup) {
	g := router.Group("recycle")
	initRunRecycleList(g)
	initRunRecycleRestore(g)
	initRunRecycleDelete(g)
	initRunRecycleClear(g)
}
