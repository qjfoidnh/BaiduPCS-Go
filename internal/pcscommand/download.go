package pcscommand

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsdownload"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/taskframework"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/downloader"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/transfer"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

type (
	//DownloadOptions 下载可选参数
	DownloadOptions struct {
		IsTest               bool
		IsPrintStatus        bool
		IsExecutedPermission bool
		IsOverwrite          bool
		DownloadMode         pcsdownload.DownloadMode
		SaveTo               string
		Parallel             int
		Load                 int
		MaxRetry             int
		NoCheck              bool
		ModifyMTime          bool
		FullPath             bool
		LinkPrefer           int
	}

	// LocateDownloadOption 获取下载链接可选参数
	LocateDownloadOption struct {
		FromPan bool
	}
)

func downloadPrintFormat(load int) string {
	if load <= 1 {
		return pcsdownload.DefaultPrintFormat
	}
	return "[%s] ↓ %s/%s %s/s in %s, left %s ...\n"
}

// RunDownload 执行下载网盘内文件
func RunDownload(paths []string, options *DownloadOptions) {
	if options == nil {
		options = &DownloadOptions{}
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
	cfg := &downloader.Config{
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

	paths, err := matchPathByShellPattern(paths...)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("\n")
	fmt.Printf("[0] 提示: 当前下载最大并发量为: %d, 下载缓存为: %d\n", options.Parallel, cfg.CacheSize)

	var (
		pcs       = GetBaiduPCS()
		loadCount = 0
	)

	// 预测要下载的文件数量
	file_dir_list := make([]*baidupcs.FileDirectory,0,10)
	for k := range paths {
		pcs.FilesDirectoriesRecurseList(paths[k], baidupcs.DefaultOrderOptions, func(depth int, _ string, fd *baidupcs.FileDirectory, pcsError pcserror.Error) bool {
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

	var (
		executor = taskframework.TaskExecutor{
			IsFailedDeque: true, // 统计失败的列表
		}
		statistic = &pcsdownload.DownloadStatistic{}
	)

	// 处理队列, 小文件优先下载
	sort.Slice(file_dir_list, func(i, j int) bool {
		return file_dir_list[i].Size < file_dir_list[j].Size
	})
	for _,v := range file_dir_list {
		newCfg := *cfg
		unit := pcsdownload.DownloadTaskUnit{
			Cfg:                  &newCfg, // 复制一份新的cfg
			PCS:                  pcs,
			VerbosePrinter:       pcsCommandVerbose,
			PrintFormat:          downloadPrintFormat(options.Load),
			ParentTaskExecutor:   &executor,
			DownloadStatistic:    statistic,
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
		executor.SetParallel(loadCount)
		// 设置储存的路径
		vPath := v.Path
		if !options.FullPath {
			vPath = filepath.Join(v.PreBase, filepath.Base(v.Path))
		}
		if options.SaveTo != "" {
			unit.SavePath = filepath.Join(options.SaveTo, vPath)
		} else {
			// 使用默认的保存路径
			unit.SavePath = GetActiveUser().GetSavePath(vPath)
		}
		info := executor.Append(&unit, options.MaxRetry)
		fmt.Printf("[%s] 加入下载队列: %s\n", info.Id(), v.Path)
	}

	// 开始计时
	statistic.StartTimer()

	// 开始执行
	executor.Execute()

	fmt.Printf("\n下载结束, 时间: %s, 数据总量: %s\n", statistic.Elapsed()/1e6*1e6, converter.ConvertFileSize(statistic.TotalSize()))

	// 输出失败的文件列表
	failedList := executor.FailedDeque()
	if failedList.Size() != 0 {
		fmt.Printf("以下文件下载失败: \n")
		tb := pcstable.NewTable(os.Stdout)
		for e := failedList.Shift(); e != nil; e = failedList.Shift() {
			item := e.(*taskframework.TaskInfoItem)
			tb.Append([]string{item.Info.Id(), item.Unit.(*pcsdownload.DownloadTaskUnit).PcsPath})
		}
		tb.Render()
	}
}
