package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
	"github.com/iikira/BaiduPCS-Go/internal/pcsfunctions/pcsupload"
	"github.com/iikira/BaiduPCS-Go/pcstable"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/iikira/BaiduPCS-Go/pcsutil/checksum"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/BaiduPCS-Go/pcsutil/taskframework"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	// DefaultUploadMaxRetry 默认上传失败最大重试次数
	DefaultUploadMaxRetry = 3
)

type (
	// UploadOptions 上传可选项
	UploadOptions struct {
		Parallel      int
		MaxRetry      int
		NoRapidUpload bool
		NoSplitFile   bool // 禁用分片上传
	}
)

// RunRapidUpload 执行秒传文件, 前提是知道文件的大小, md5, 前256KB切片的 md5, crc32
func RunRapidUpload(targetPath, contentMD5, sliceMD5, crc32 string, length int64) {
	err := matchPathByShellPatternOnce(&targetPath)
	if err != nil {
		fmt.Printf("警告: %s, 获取网盘路径 %s 错误, %s\n", baidupcs.OperationRapidUpload, targetPath, err)
	}

	err = GetBaiduPCS().RapidUpload(targetPath, contentMD5, sliceMD5, crc32, length)
	if err != nil {
		fmt.Printf("%s失败, 消息: %s\n", baidupcs.OperationRapidUpload, err)
		return
	}

	fmt.Printf("%s成功, 保存到网盘路径: %s\n", baidupcs.OperationRapidUpload, targetPath)
	return
}

// RunCreateSuperFile 执行分片上传—合并分片文件
func RunCreateSuperFile(targetPath string, blockList ...string) {
	err := matchPathByShellPatternOnce(&targetPath)
	if err != nil {
		fmt.Printf("警告: %s, 获取网盘路径 %s 错误, %s\n", baidupcs.OperationUploadCreateSuperFile, targetPath, err)
	}

	err = GetBaiduPCS().UploadCreateSuperFile(true, targetPath, blockList...)
	if err != nil {
		fmt.Printf("%s失败, 消息: %s\n", baidupcs.OperationUploadCreateSuperFile, err)
		return
	}

	fmt.Printf("%s成功, 保存到网盘路径: %s\n", baidupcs.OperationUploadCreateSuperFile, targetPath)
	return
}

// RunUpload 执行文件上传
func RunUpload(localPaths []string, savePath string, opt *UploadOptions) {
	if opt == nil {
		opt = &UploadOptions{}
	}

	// 检测opt
	if opt.Parallel <= 0 {
		opt.Parallel = pcsconfig.Config.MaxUploadParallel
	}

	if opt.MaxRetry < 0 {
		opt.MaxRetry = DefaultUploadMaxRetry
	}

	err := matchPathByShellPatternOnce(&savePath)
	if err != nil {
		fmt.Printf("警告: 上传文件, 获取网盘路径 %s 错误, %s\n", savePath, err)
	}

	switch len(localPaths) {
	case 0:
		fmt.Printf("本地路径为空\n")
		return
	}

	// 打开上传状态
	uploadDatabase, err := pcsupload.NewUploadingDatabase()
	if err != nil {
		fmt.Printf("打开上传未完成数据库错误: %s\n", err)
		return
	}
	defer uploadDatabase.Close()

	var (
		pcs = GetBaiduPCS()
		// 使用 task framework
		executor = &taskframework.TaskExecutor{
			IsFailedDeque: true, // 失败统计
		}
		subSavePath string
		// 统计
		statistic = &pcsupload.UploadStatistic{}
	)

	statistic.StartTimer() // 开始计时

	for k := range localPaths {
		walkedFiles, err := pcsutil.WalkDir(localPaths[k], "")
		if err != nil {
			fmt.Printf("警告: 遍历错误: %s\n", err)
			continue
		}

		for k3 := range walkedFiles {
			var localPathDir string
			// 针对 windows 的目录处理
			if os.PathSeparator == '\\' {
				walkedFiles[k3] = pcsutil.ConvertToUnixPathSeparator(walkedFiles[k3])
				localPathDir = pcsutil.ConvertToUnixPathSeparator(filepath.Dir(localPaths[k]))
			} else {
				localPathDir = filepath.Dir(localPaths[k])
			}

			// 避免去除文件名开头的"."
			if localPathDir == "." {
				localPathDir = ""
			}

			subSavePath = strings.TrimPrefix(walkedFiles[k3], localPathDir)

			info := executor.Append(&pcsupload.UploadTaskUnit{
				LocalFileChecksum: checksum.NewLocalFileChecksum(walkedFiles[k3], int(baidupcs.SliceMD5Size)),
				SavePath:          path.Clean(savePath + baidupcs.PathSeparator + subSavePath),
				PCS:               pcs,
				UploadingDatabase: uploadDatabase,
				Parallel:          opt.Parallel,
				NoRapidUpload:     opt.NoRapidUpload,
				NoSplitFile:       opt.NoSplitFile,
				UploadStatistic:   statistic,
			}, opt.MaxRetry)
			fmt.Printf("[%s] 加入上传队列: %s\n", info.Id(), walkedFiles[k3])
		}
	}

	// 没有添加任何任务
	if executor.Count() == 0 {
		fmt.Printf("未检测到上传的文件.\n")
		return
	}

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
}
