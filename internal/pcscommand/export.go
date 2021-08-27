package pcscommand

import (
	"container/list"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/pcstime"
)

type (
	etask struct {
		*ListTask
		path     string
		rootPath string
		fd       *baidupcs.FileDirectory
		err      pcserror.Error
	}

	// ExportOptions 导出可选项
	ExportOptions struct {
		RootPath   string // 根路径
		SavePath   string // 输出路径
		MaxRetry   int
		Recursive  bool
		LinkFormat bool
		StdOut     bool
	}
)

func (task *etask) handleExportTaskError(l *list.List, failedList *list.List) {
	if task.err == nil {
		return
	}

	// 不重试
	switch task.err.GetError() {
	case baidupcs.ErrGetRapidUploadInfoMD5NotFound, baidupcs.ErrGetRapidUploadInfoCrc32NotFound:
		fmt.Printf("[%d] - [%s] 导出失败, 可能是服务器未刷新文件的md5, 请过一段时间再试一试\n", task.ID, task.path)
		failedList.PushBack(task)
		return
	case baidupcs.ErrFileTooLarge:
		fmt.Printf("[%d] - [%s] 导出失败, 文件大于20GB, 无法导出\n", task.ID, task.path)
		failedList.PushBack(task)
		return
	}

	// 未达到失败重试最大次数, 将任务推送到队列末尾
	if task.retry < task.MaxRetry {
		task.retry++
		fmt.Printf("[%d] - [%s] 导出错误, %s, 重试 %d/%d\n", task.ID, task.path, task.err, task.retry, task.MaxRetry)
		l.PushBack(task)
		time.Sleep(3 * time.Duration(task.retry) * time.Second)
	} else {
		fmt.Printf("[%d] - [%s] 导出错误, %s\n", task.ID, task.path, task.err)
		failedList.PushBack(task)
	}
}

func changeRootPath(dstRootPath, dstPath, srcRootPath string) string {
	if srcRootPath == "" {
		return dstPath
	}
	return path.Join(srcRootPath, strings.TrimPrefix(dstPath, dstRootPath))
}

// GetExportFilename 获取导出路径
func GetExportFilename() string {
	return "BaiduPCS-Go_export_" + pcstime.BeijingTimeOption("") + ".txt"
}

// RunExport 执行导出文件和目录
func RunExport(pcspaths []string, opt *ExportOptions) {
	if opt == nil {
		opt = &ExportOptions{}
	}

	if opt.SavePath == "" {
		opt.SavePath = GetExportFilename()
	}

	pcspaths, err := matchPathByShellPattern(pcspaths...)
	if err != nil {
		fmt.Println(err)
		return
	}
	saveFile, err := os.OpenFile(opt.SavePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil { // 不可写
		if !opt.StdOut {
			fmt.Printf("%s\n", err)
			return
		}
	}
	defer saveFile.Close()
	if !opt.StdOut {
		fmt.Printf("导出的信息将保存在: %s\n", opt.SavePath)
	}

	var (
		au         = GetActiveUser()
		pcs        = GetBaiduPCS()
		l          = list.New()
		failedList = list.New()
		writeErr   error
		id         int
	)

	for id = range pcspaths {
		var rootPath string
		if pcspaths[id] == au.Workdir {
			rootPath = pcspaths[id]
		} else {
			rootPath = path.Dir(pcspaths[id])
		}
		// 加入队列
		l.PushBack(&etask{
			ListTask: &ListTask{
				ID:       id,
				MaxRetry: opt.MaxRetry,
			},
			path:     pcspaths[id],
			rootPath: rootPath,
		})
	}

	for {
		e := l.Front()
		if e == nil { // 结束
			break
		}

		l.Remove(e) // 载入任务后, 移除队列

		task := e.Value.(*etask)
		root := task.fd == nil

		// 获取文件信息
		if task.fd == nil { // 第一次初始化
			fd, pcsError := pcs.FilesDirectoriesMeta(task.path)
			if pcsError != nil {
				task.err = pcsError
				task.handleExportTaskError(l, failedList)
				continue
			}
			task.fd = fd
		}

		if task.fd.Isdir { // 导出目录
			if !root && !opt.Recursive { // 非递归
				continue
			}

			fds, pcsError := pcs.FilesDirectoriesList(task.path, baidupcs.DefaultOrderOptions)
			if pcsError != nil {
				task.err = pcsError
				task.handleExportTaskError(l, failedList)
				continue
			}

			if len(fds) == 0 && !opt.StdOut {
				_, writeErr = saveFile.Write(converter.ToBytes(fmt.Sprintf("BaiduPCS-Go mkdir \"%s\"\n", changeRootPath(task.rootPath, task.path, opt.RootPath))))
				if writeErr != nil {
					fmt.Printf("写入文件失败: %s\n", writeErr)
					return // 直接返回
				}
				fmt.Printf("[%d] - [%s] 导出成功\n", task.ID, task.path)
				continue
			}

			// 加入队列
			for _, fd := range fds {
				// 加入队列
				id++
				l.PushBack(&etask{
					ListTask: &ListTask{
						ID:       id,
						MaxRetry: opt.MaxRetry,
					},
					path:     fd.Path,
					fd:       fd,
					rootPath: task.rootPath,
				})
			}
			continue
		}

		rinfo, pcsError := pcs.ExportByFileInfo(task.fd)
		if pcsError != nil {
			task.err = pcsError
			task.handleExportTaskError(l, failedList)
			continue
		}
		var outTemplate = fmt.Sprintf("BaiduPCS-Go rapidupload -length=%d -md5=%s -slicemd5=%s -crc32=%s \"%s\"\n", rinfo.ContentLength, rinfo.ContentMD5, rinfo.SliceMD5, rinfo.ContentCrc32, changeRootPath(task.rootPath, task.path, opt.RootPath))
		if opt.LinkFormat {
			outTemplate = fmt.Sprintf("%s#%s#%d#%s\n", rinfo.ContentMD5, rinfo.SliceMD5, rinfo.ContentLength, path.Base(task.path))
		}
		if opt.StdOut {
			fmt.Print(outTemplate)
		} else {
			_, writeErr = saveFile.Write(converter.ToBytes(outTemplate))
			if writeErr != nil {
				fmt.Printf("写入文件失败: %s\n", writeErr)
				return // 直接返回
			}

			fmt.Printf("[%d] - [%s] 导出成功\n", task.ID, task.path)
		}
	}
	if opt.StdOut {
		os.Remove(opt.SavePath)
		fmt.Println("导出完毕")
	}

	if failedList.Len() > 0 {
		fmt.Printf("\n以下目录导出失败: \n")
		fmt.Printf("%s\n", strings.Repeat("-", 100))
		for e := failedList.Front(); e != nil; e = e.Next() {
			et := e.Value.(*etask)
			fmt.Printf("[%d] %s\n", et.ID, et.path)
		}
	}
}
