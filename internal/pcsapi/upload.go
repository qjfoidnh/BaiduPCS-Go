package pcsapi

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsupload"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/taskframework"
)

type UploadStructure struct {
	UploadThread []string `json:"upload_thread" form:"upload_thread"`
	Retry        int      `json:"retry" form:"retry"`
	Parallel     int      `json:"parallel" form:"parallel"`
	Norapid      bool     `json:"norapid,omitempty" form:"norapid,omitempty"`
	NoSplit      bool     `json:"no_split,omitempty" form:"no_split,omitempty"`
	Policy       string   `json:"policy" form:"policy"`
	LocalPaths   []string `json:"local_paths" form:"local_paths" binding:"required"`
	SavePath     string   `json:"save_path" form:"save_path"`
}

func parseUploadArgs(args *UploadStructure) (opt *pcscommand.UploadOptions, err error) {
	opt = &pcscommand.UploadOptions{
		Parallel:      args.Parallel,
		MaxRetry:      args.Retry,
		Load:          args.Parallel,
		NoRapidUpload: args.Norapid,
		NoSplitFile:   args.NoSplit,
		Policy:        args.Policy,
	}
	// 检测opt
	if opt.Parallel <= 0 {
		opt.Parallel = pcsconfig.Config.MaxUploadParallel
	}
	opt.NoFilenameCheck = pcsconfig.Config.IgnoreIllegal

	if opt.MaxRetry < 0 {
		opt.MaxRetry = pcscommand.DefaultUploadMaxRetry
	}

	if opt.Load <= 0 {
		opt.Load = pcsconfig.Config.MaxUploadLoad
	}

	if opt.Policy != "fail" && opt.Policy != "newcopy" && opt.Policy != "overwrite" && opt.Policy != "skip" && opt.Policy != "rsync" {
		opt.Policy = pcsconfig.Config.UPolicy
	}
	err = matchPathByShellPatternOnce(&args.SavePath)
	if err != nil {
		return
	}
	if len(args.LocalPaths) <= 0 {
		err = fmt.Errorf("本地路径为空")
	}
	return
}

func uploadPrintFormat(load int) string {
	if load <= 1 {
		return pcsupload.DefaultPrintFormat
	}
	return "[%s] ↑ %s/%s %s/s in %s ...\n"
}

func runUpload(ctx *gin.Context) {
	args := UploadStructure{
		Policy: "",
	}
	// 解析参数
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("upload command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	opt, err := parseUploadArgs(&args)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 打开上传状态
	uploadDatabase, err := pcsupload.NewUploadingDatabase()
	if err != nil {
		fmt.Printf("打开上传未完成数据库错误: %s\n", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}
	defer uploadDatabase.Close()

	var (
		pcs = pcscommand.GetBaiduPCS()
		// 使用 task framework
		executor = &taskframework.TaskExecutor{
			IsFailedDeque: true, // 失败统计
		}
		subSavePath string
		// 统计
		statistic = &pcsupload.UploadStatistic{}
	)
	fmt.Print("\n")
	fmt.Printf("[0] 提示: 当前上传单个文件最大并发量为: %d, 最大同时上传文件数为: %d\n", opt.Parallel, opt.Load)

	statistic.StartTimer() // 开始计时

	LoadCount := 0
	var errs []error

	for k := range args.LocalPaths {
		walkedFiles, err := pcsutil.WalkDir(args.LocalPaths[k], "")
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for k3 := range walkedFiles {
			var localPathDir string
			// 针对 windows的目录处理
			if os.PathSeparator == '\\' {
				walkedFiles[k3] = pcsutil.ConvertToUnixPathSeparator(walkedFiles[k3])
				localPathDir = pcsutil.ConvertToUnixPathSeparator(filepath.Dir(args.LocalPaths[k]))
			} else {
				localPathDir = filepath.Dir(args.LocalPaths[k])
			}
			// 避免去除文件开头的"."
			if localPathDir == "." {
				localPathDir = ""
			}
			if len(args.LocalPaths) == 1 && len(walkedFiles) == 1 {
				opt.Load = 1
			}
			subSavePath = strings.TrimPrefix(walkedFiles[k3], localPathDir)
			if !opt.NoFilenameCheck && !pcsutil.ChPathLegal(walkedFiles[k3]) {
				errs = append(errs, fmt.Errorf("[0] %s 文件路径含有非法字符，已跳过", walkedFiles[k3]))
				continue
			}
			LoadCount++
			info := executor.Append(&pcsupload.UploadTaskUnit{
				LocalFileChecksum: checksum.NewLocalFileChecksum(walkedFiles[k3], int(baidupcs.SliceMD5Size)),
				SavePath:          path.Clean(args.SavePath + baidupcs.PathSeparator + subSavePath),
				PCS:               pcs,
				UploadingDatabase: uploadDatabase,
				Parallel:          opt.Parallel,
				PrintFormat:       uploadPrintFormat(opt.Load),
				NoRapidUpload:     opt.NoRapidUpload,
				NoSplitFile:       opt.NoSplitFile,
				UploadStatistic:   statistic,
				Policy:            opt.Policy,
			}, opt.MaxRetry)
			if LoadCount >= opt.Load {
				LoadCount = opt.Load
			}
			fmt.Printf("[%s] 加入上传队列: %s\n", info.Id(), walkedFiles[k3])
		}
	}

	if len(errs) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"error": errs,
		})
		return
	}

	// 没有添加任何任务
	if executor.Count() <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"error": "未检测到上传的文件",
		})
		return
	}

	// 设置上传文件并发数
	executor.SetParallel(LoadCount)

	res := gin.H{
		"message": fmt.Sprintf("共开始%d个上传任务", LoadCount),
	}
	// 没有添加任何任务
	if executor.Count() <= 0 {
		res["error"] = "未检测到上传的文件"
	} else if len(errs) > 0 {
		// 添加上传任务失败
		res["error"] = errs
	}

	ctx.JSON(http.StatusOK, res)

	// 后台启动协程开始上传
	go func() {
		// 执行上传任务
		executor.Execute()

		fmt.Printf("\n")
		fmt.Printf("上传结束, 时间: %s, 总大小: %s\n", statistic.Elapsed()/1e6*1e6, converter.ConvertFileSize(statistic.TotalSize()))

		// 输出上传失败的文件列表
		failedList := executor.FailedDeque()
		if failedList.Size() != 0 {
			fmt.Printf("以下文件上传失败: \n")
			tb := pcstable.NewTable(os.Stdout)
			for e := failedList.Shift(); e != nil; e = failedList.Shift() {
				item := e.(*taskframework.TaskInfoItem)
				tb.Append([]string{item.Info.Id(), item.Unit.(*pcsupload.UploadTaskUnit).LocalFileChecksum.Path})
			}
			tb.Render()
		}
	}()
}

// RunRapidUpload 执行秒传文件, 前提是知道文件的大小, md5, 前256KB切片的 md5, crc32
func RunRapidUpload(targetPath, contentMD5, sliceMD5, crc32 string, length int64) (err error) {
	dirname := path.Dir(targetPath)
	err = matchPathByShellPatternOnce(&dirname)
	if err != nil {
		fmt.Printf("警告: %s, 获取网盘路径 %s 错误, %s\n", baidupcs.OperationRapidUpload, dirname, err)
	}
	err = pcscommand.GetBaiduPCS().RapidUpload(targetPath, contentMD5, sliceMD5, crc32, length)
	if err != nil {
		fmt.Printf("%s失败, 消息: %s\n", baidupcs.OperationRapidUpload, err)
		return
	}

	fmt.Printf("%s成功, 保存到网盘路径: %s\n", baidupcs.OperationRapidUpload, targetPath)
	return
}

func initRunUpload(group *gin.RouterGroup) {
	group.POST("upload", runUpload)
}
