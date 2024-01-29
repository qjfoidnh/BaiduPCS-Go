package pcsapi

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsdownload"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/taskframework"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/downloader"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/transfer"
)

type DownloadStructure struct {
	// 将下载的文件直接保存到当前工作目录
	Save bool `json:"save,omitempty" form:"save,omitempty"`
	// 将下载的文件直接保存到指定的目录
	SaveTo string `json:"save_to" form:"save_to"`
	// 下载模式, 可选值: pcs, stream, locate, 默认为 locate,
	Mode string `json:"mode" form:"mode"`
	// 测试下载, 此操作不会保存文件到本地
	IsTest bool `json:"is_test,omitempty" form:"is_test,omitempty"`
	// 是否输出所有线程的工作状态,默认为是
	IsPrintStatus bool `json:"is_printStatus,omitempty" form:"is_printStatus,omitempty"`
	// 为文件加上执行权限, (windows系统无效)
	IsExecutedPermission bool `json:"is_executed_permission,omitempty" form:"is_executed_permission,omitempty"`
	// overwrite, 覆盖已存在的文件
	IsOverwrite bool `json:"is_overwrite,omitempty" form:"is_overwrite,omitempty"`
	// 下载线程数
	Parallel int `json:"parallel,omitempty" form:"parallel,omitempty"`
	// 指定同时进行下载文件的数量
	Load int `json:"load,omitempty" form:"load,omitempty"`
	// 下载失败最大重试次数
	MaxRetry int `json:"max_retry,omitempty" form:"max_retry,omitempty"`
	// 下载文件完成后不校验文件
	NoCheck bool `json:"no_check,omitempty" form:"no_check,omitempty"`
	// 使用备选下载链接中的第几个，默认第一个
	LinkPrefer int `json:"link_prefer,omitempty" form:"link_prefer,omitempty"`
	// 将本地文件的修改时间设置为服务器上的修改时间
	ModifyMTime bool `json:"modifyMTime,omitempty" form:"modifyMTime,omitempty"`
	// 以网盘完整路径保存到本地
	FullPath bool `json:"full_path,omitempty" form:"full_path,omitempty"`
	// 下载路径
	Paths []string `json:"paths" form:"paths" binding:"required"`
}

// 全局变量，包含执行器，下载事件管道
var (
	Executor = taskframework.TaskExecutor{
		IsFailedDeque: true, // 统计失败的列表
	}
	Statistic = &pcsdownload.DownloadStatistic{}
	// 下载事件通道
	dl_channel = make(chan sse.Event, runtime.NumCPU())
)

// 处理下载设置
func dealArgs(args *DownloadStructure) (config *downloader.Config, options *pcscommand.DownloadOptions, err error) {

	// 处理saveTo
	var (
		saveTo string
	)
	if args.Save {
		saveTo = "."
	} else if args.SaveTo != "" {
		saveTo = filepath.Clean(args.SaveTo)
	}
	var (
		downloadMode pcsdownload.DownloadMode
	)
	// 处理解析downloadMode
	switch args.Mode {
	case "pcs":
		downloadMode = pcsdownload.DownloadModePCS
	case "stream":
		downloadMode = pcsdownload.DownloadModeStreaming
	case "locate":
		downloadMode = pcsdownload.DownloadModeLocate
	default:
		err = fmt.Errorf("下载方式解析失败")
		return
	}
	options = &pcscommand.DownloadOptions{
		IsTest:               args.IsTest,
		IsPrintStatus:        args.IsPrintStatus,
		IsExecutedPermission: args.IsExecutedPermission,
		IsOverwrite:          args.IsOverwrite,
		DownloadMode:         downloadMode,
		SaveTo:               saveTo,
		Parallel:             args.Parallel,
		Load:                 args.Load,
		MaxRetry:             args.MaxRetry,
		NoCheck:              args.NoCheck,
		LinkPrefer:           args.LinkPrefer,
		ModifyMTime:          args.ModifyMTime,
		FullPath:             args.FullPath,
	}

	if options.Load <= 0 {
		options.Load = pcsconfig.Config.MaxDownloadLoad
	}

	if options.MaxRetry < 0 {
		options.MaxRetry = pcsdownload.DefaultDownloadMaxRetry
	}

	if !options.NoCheck {
		options.NoCheck = pcsconfig.Config.NoCheck
	}

	if runtime.GOOS == "windows" {
		// windows下不加执行权限
		options.IsExecutedPermission = false
	}
	// 设置下载配置
	config = &downloader.Config{
		Mode:                       transfer.RangeGenMode_BlockSize,
		CacheSize:                  pcsconfig.Config.CacheSize,
		BlockSize:                  baidupcs.InitRangeSize,
		MaxRate:                    pcsconfig.Config.MaxDownloadRate,
		InstanceStateStorageFormat: downloader.InstanceStateStorageFormatProto3,
		IsTest:                     options.IsTest,
		TryHTTP:                    !pcsconfig.Config.EnableHTTPS,
	}

	// 设置下载最大并发量
	if options.Parallel < 1 {
		options.Parallel = pcsconfig.Config.MaxParallel
	}

	return
}

// 下载文件
func runDownload(ctx *gin.Context) {
	// 设置默认值
	args := DownloadStructure{
		Parallel:   0,
		Mode:       "locate",
		Load:       1,
		MaxRetry:   pcsdownload.DefaultDownloadMaxRetry,
		LinkPrefer: 1,
	}
	// dl_channel <- sse.Event{
	// 	Event: "test",
	// 	Data:  args,
	// }
	// 解析参数
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("download command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 解析设置并生成config

	var (
		cfg     *downloader.Config
		options *pcscommand.DownloadOptions
		err     error
	)

	cfg, options, err = dealArgs(&args)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	paths, err := matchPathByShellPattern(args.Paths...)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	var (
		pcs       = pcscommand.GetBaiduPCS()
		loadCount = 0
	)

	file_dir_list := make([]*baidupcs.FileDirectory, 0, 10)
	for k := range paths {
		pcs.FilesDirectoriesRecurseList(paths[k], baidupcs.DefaultOrderOptions, func(depth int, fdPath string, fd *baidupcs.FileDirectory, pcsError pcserror.Error) bool {
			if pcsError != nil {
				pcsCommandVerbose.Warnf("%s\n", pcsError)
				return true
			}
			file_dir_list = append(file_dir_list, fd)
			// 忽略统计文件夹数量
			if !fd.Isdir {
				loadCount++
				if loadCount >= options.Load {
					loadCount = options.Load
				}
			}
			return true
		})
	}

	// 修改Load, 设置MaxParallel
	if loadCount > 0 {
		options.Load = loadCount
		// 取平均值
		cfg.MaxParallel = pcsconfig.AverageParallel(options.Parallel, loadCount)
	} else {
		cfg.MaxParallel = options.Parallel
	}

	// 处理队列, 小文件优先下载
	sort.Slice(file_dir_list, func(i, j int) bool {
		return file_dir_list[i].Size < file_dir_list[j].Size
	})
	for _, v := range file_dir_list {
		newCfg := *cfg
		unit := pcsdownload.DownloadTaskUnit{
			Cfg:            &newCfg, //复制一份新的cfg
			PCS:            pcs,
			VerbosePrinter: pcsCommandVerbose,
			PrintFormat: func(load int) string {
				if load <= 1 {
					return pcsdownload.DefaultPrintFormat
				}
				return "[%s] ↓ %s/%s %s/s in %s, left %s ...\n"
			}(options.Load),
			ParentTaskExecutor:   &Executor,
			DownloadStatistic:    Statistic,
			IsPrintStatus:        options.IsPrintStatus,
			IsExecutedPermission: options.IsExecutedPermission,
			IsOverwrite:          options.IsOverwrite,
			NoCheck:              options.NoCheck,
			DlinkPrefer:          options.LinkPrefer,
			DownloadMode:         options.DownloadMode,
			ModifyMTime:          options.ModifyMTime,
			PcsPath:              v.Path,
			FileInfo:             v,
		}
		// 设置下载并发数
		Executor.SetParallel(loadCount)
		// 设置储存的路径
		vPath := v.Path
		if !options.FullPath {
			vPath = filepath.Join(v.PreBase, filepath.Base(v.Path))
		}
		if options.SaveTo != "" {
			unit.SavePath = filepath.Join(options.SaveTo, vPath)
		} else {
			// 使用默认的保存路径
			unit.SavePath = pcscommand.GetActiveUser().GetSavePath(vPath)
		}
		info := Executor.Append(&unit, options.MaxRetry)
		fmt.Printf("[%s] 加入下载队列: %s\n", info.Id(), v.Path)
	}
	//启动协程，开始下载
	go func() {
		Executor.Execute()
	}()
}

// 将runDownload挂载到路由列表
func initRunDownload(group *gin.RouterGroup) {
	// group.GET("download", runDownload)
	group.POST("download", runDownload)
}
