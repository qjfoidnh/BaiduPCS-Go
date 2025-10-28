package pcsupload

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/taskframework"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/rio"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/uploader"
	"path"
	"strings"
	"time"
)

type (
	// StepUpload 上传步骤
	StepUpload int

	// UploadTaskUnit 上传的任务单元
	UploadTaskUnit struct {
		LocalFileChecksum *checksum.LocalFileChecksum // 要上传的本地文件详情
		Step              StepUpload
		SavePath          string // 保存路径
		PrintFormat       string

		PCS               *baidupcs.BaiduPCS
		UploadingDatabase *UploadingDatabase // 数据库
		Parallel          int
		NoRapidUpload     bool   // 禁用秒传
		NoSplitFile       bool   // 禁用分片上传
		Policy            string // 上传重名文件策略

		UploadStatistic *UploadStatistic

		taskInfo *taskframework.TaskInfo
		panDir   string
		panFile  string
		state    *uploader.InstanceState
	}
)

const (
	// StepUploadInit 初始化步骤
	StepUploadInit StepUpload = iota
	// StepUploadRapidUpload 秒传步骤
	StepUploadRapidUpload
	// StepUploadUpload 正常上传步骤
	StepUploadUpload
	// JustGoon 跳过上传
	JustGoon
)

const (
	StrUploadFailed    = "上传文件失败"
	DefaultPrintFormat = "\r[%s] ↑ %s/%s %s/s in %s ............"
	DefaultContentSize = 4 * converter.KB
)

func (utu *UploadTaskUnit) SetTaskInfo(taskInfo *taskframework.TaskInfo) {
	utu.taskInfo = taskInfo
}

// prepareFile 解析文件阶段
func (utu *UploadTaskUnit) prepareFile() {
	// 解析文件保存路径
	var (
		panDir, panFile = path.Split(utu.SavePath)
	)
	utu.panDir = path.Clean(panDir)
	utu.panFile = panFile

	// 检测断点续传
	// 2025.10.26 不再支持续传, 关闭检测
	//utu.state = utu.UploadingDatabase.Search(&utu.LocalFileChecksum.LocalFileMeta)
	//if utu.state != nil || utu.LocalFileChecksum.LocalFileMeta.BlocksList != nil { // 读取到了md5分片信息
	//	utu.Step = StepUploadRapidUpload
	//	fmt.Printf("[%s] 检测到断点信息, 准备续传...\n", utu.taskInfo.Id())
	//	return
	//}
	utu.state = &uploader.InstanceState{}

	if utu.LocalFileChecksum.Length >= baidupcs.RecommendedUploadSize {
		fmt.Printf("[%s] 文件超过32GB, 上传有可能失败, 建议分割文件...\n", utu.taskInfo.Id())
	}

	if utu.LocalFileChecksum.Length > baidupcs.MinCheckLeftSpaceThreshold {
		freeSpace, err := utu.PCS.SpaceLeftInfo()
		if err == nil && freeSpace < utu.LocalFileChecksum.Length {
			fmt.Printf("[%s] 目标文件大小超过剩余空间, 跳过...\n", utu.taskInfo.Id())
			utu.Step = JustGoon
			return
		}
	}

	if utu.NoRapidUpload {
		//fmt.Printf("[%s] 注意: 跳过秒传将无法使用断点续传...\n", utu.taskInfo.Id())
		pcsError, jsonData := utu.PCS.FakeRapidUpload(utu.SavePath, utu.Policy, utu.LocalFileChecksum.Length)
		if pcsError != nil {
			errcode := pcsError.GetRemoteErrCode()
			if errcode != 114514 && errcode != 1919810 {
				fmt.Printf("[%s] 跳过秒传失败, 开始秒传...\n", utu.taskInfo.Id())
				utu.Step = StepUploadRapidUpload
				return
			} else {
				fmt.Printf("[%s] 目标文件已存在, 跳过...\n", utu.taskInfo.Id())
				utu.Step = JustGoon
				return
			}
		}
		utu.state = &uploader.InstanceState{
			Uploadid: jsonData.UploadID,
		}
		utu.UploadingDatabase.Save()
		utu.Step = StepUploadUpload
		return
	}

	// 秒传不分文件大小一律进行
	utu.Step = StepUploadRapidUpload
}

// rapidUpload 执行秒传
func (utu *UploadTaskUnit) rapidUpload() (isContinue bool, result *taskframework.TaskUnitRunResult) {
	utu.Step = StepUploadRapidUpload

	// TODO: 建立一个通过百度错误码判断重试的函数
	result = &taskframework.TaskUnitRunResult{}

	fdl, pcsError := utu.PCS.CacheFilesDirectoriesList(utu.panDir, baidupcs.DefaultOrderOptions)
	if pcsError != nil {
		switch pcsError.GetErrType() {
		case pcserror.ErrTypeRemoteError:
			switch pcsError.GetRemoteErrCode() {
			case 31066:
			// file does not exist
			// 不缓存文件夹
			default:
				// 其他百度服务器错误, 不重试
				result.ResultMessage = "获取文件列表错误"
				result.Err = pcsError
				return
			}
		default:
			// 未知错误, 重试
			result.ResultMessage = "获取文件列表错误"
			result.NeedRetry = true
			result.Err = pcsError
			return
		}
	}

	fmt.Printf("[%s] 开始计算文件元信息, 请稍候...\n", utu.taskInfo.Id())

	// 经测试, 文件的 crc32 值并非秒传文件所必需
	if utu.LocalFileChecksum.LocalFileMeta.MD5 == nil || utu.LocalFileChecksum.LocalFileMeta.SliceMD5 == nil {
		err := utu.LocalFileChecksum.Sum(checksum.CHECKSUM_MD5 | checksum.CHECKSUM_SLICE_MD5)
		if err != nil {
			// 不重试
			result.ResultMessage = "计算文件秒传信息错误"
			result.Err = err
			return
		}
	}

	// 检测缓存, 通过文件的md5值判断本地文件和网盘文件是否一样
	if fdl != nil {
		for _, fd := range fdl {
			if fd.Filename == utu.panFile {
				decodedMD5, _ := hex.DecodeString(fd.MD5)
				// TODO: fd.MD5 有可能是错误的
				if (utu.Policy == baidupcs.SkipPolicy) || (bytes.Compare(decodedMD5, utu.LocalFileChecksum.MD5) == 0) {
					fmt.Printf("[%s] 目标文件, %s, 已存在, 跳过...\n", utu.taskInfo.Id(), utu.SavePath)
					result.Succeed = true // 成功
					return
				}
			}
		}
	}

	uk, pcsError := utu.PCS.CacheUK()
	if pcsError != nil {
		result.ResultMessage = "获取用户uk错误, 请确保登录信息包含了STOKEN"
		result.Err = pcsError
		return
	}
	currentTime := time.Now().Unix()
	offset, err := creaetDataOffset(hex.EncodeToString(utu.LocalFileChecksum.MD5), uk, currentTime, utu.LocalFileChecksum.Length, DefaultContentSize)
	if err != nil {
		result.ResultMessage = "计算文件偏移量错误"
		result.Err = err
		return
	}
	dataContent, dataLength, err := utu.LocalFileChecksum.GetSliceDataContent(offset, DefaultContentSize)
	if err != nil {
		result.ResultMessage = "读取随机文件子片段错误"
		result.Err = err
		return
	}
	b64Content := strings.TrimRight(base64.StdEncoding.EncodeToString(dataContent), "=")

	blockSize := getBlockSize(utu.LocalFileChecksum.Length)

	fmt.Printf("[%s] 开始计算文件分块md5, 请稍候...\n", utu.taskInfo.Id())
	if utu.LocalFileChecksum.LocalFileMeta.BlocksList == nil || len(utu.LocalFileChecksum.LocalFileMeta.BlocksList) == 0 {
		err = utu.LocalFileChecksum.CalculateChunkedSum(blockSize)
		if err != nil {
			// 不重试
			result.ResultMessage = "计算文件分块md5出错"
			result.Err = err
			return
		}
	}

	pcsError, jsonData := utu.PCS.RapidUpload(utu.SavePath, utu.Policy, utu.state.Uploadid, hex.EncodeToString(utu.LocalFileChecksum.MD5),
		hex.EncodeToString(utu.LocalFileChecksum.SliceMD5), b64Content, fmt.Sprint(utu.LocalFileChecksum.CRC32),
		offset, dataLength, utu.LocalFileChecksum.Length, currentTime, utu.LocalFileChecksum.BlocksList)
	if pcsError == nil {
		if jsonData.ReturnType == 2 {
			fmt.Printf("[%s] 秒传成功, 保存到网盘路径: %s\n\n", utu.taskInfo.Id(), utu.SavePath)
			// 统计
			utu.UploadStatistic.AddTotalSize(utu.LocalFileChecksum.Length)
			result.Succeed = true // 成功
			return
		}
	} else {
		// 判断配额是否已满
		switch pcsError.GetErrType() {
		// 远程服务器错误
		case pcserror.ErrTypeRemoteError:
			switch pcsError.GetRemoteErrCode() {
			case 31112: //exceed quota
				result.ResultMessage = "秒传失败, 超出配额, 网盘容量已满"
			case 114514:
				// 自定义错误码, 仅在skip策略下出现
				result.ResultMessage = StrUploadFailed
				result.Err = pcsError
				if utu.Policy == baidupcs.SkipPolicy {
					result.Extra = baidupcs.SkipPolicy
					result.Err = nil
					result.ResultMessage = fmt.Sprintf("%s 目标已存在, 跳过", utu.SavePath)
				}
				result.NeedRetry = false
				return
			case 1919810:
				// 自定义错误码, 仅在rsync策略下出现
				result.Extra = baidupcs.RsyncPolicy
				result.Err = nil
				result.ResultMessage = fmt.Sprintf("%s 目标大小未发生改变, 跳过", utu.SavePath)
				result.NeedRetry = false
				return
			}
		}
		return
	}

	fmt.Printf("[%s] 开始上传文件...\n\n", utu.taskInfo.Id())

	// 保存秒传信息
	if utu.state.Uploadid == "" {
		utu.state.Uploadid = jsonData.UploadID
	} else {
		utu.UploadingDatabase.UpdateFullBlock(&utu.LocalFileChecksum.LocalFileMeta, utu.state)
	}

	utu.UploadingDatabase.UpdateUploading(&utu.LocalFileChecksum.LocalFileMeta, utu.state)
	utu.UploadingDatabase.Save()
	isContinue = true
	return
}

// upload 上传文件
func (utu *UploadTaskUnit) upload() (result *taskframework.TaskUnitRunResult) {
	utu.Step = StepUploadUpload

	blockSize := getBlockSize(utu.LocalFileChecksum.Length)

	muer := uploader.NewMultiUploader(NewPCSUpload(utu.PCS, utu.SavePath), rio.NewFileReaderAtLen64(utu.LocalFileChecksum.GetFile()), &uploader.MultiUploaderConfig{
		Parallel:  utu.Parallel,
		BlockSize: blockSize,
		MaxRate:   pcsconfig.Config.MaxUploadRate,
		Policy:    utu.Policy,
	}, utu.SavePath)

	// 设置断点续传
	if utu.state != nil {
		muer.SetInstanceState(utu.state)
	}
	muer.OnUploadStatusEvent(func(status uploader.Status, updateChan <-chan struct{}) {
		select {
		case <-updateChan:
			if utu.state.Uploadid != "" {
				utu.UploadingDatabase.UpdateUploading(&utu.LocalFileChecksum.LocalFileMeta, muer.InstanceState())
				utu.UploadingDatabase.Save()
			}
		default:
		}

		fmt.Printf(utu.PrintFormat, utu.taskInfo.Id(),
			converter.ConvertFileSize(status.Uploaded(), 2),
			converter.ConvertFileSize(status.TotalSize(), 2),
			converter.ConvertFileSize(status.SpeedsPerSecond(), 2),
			status.TimeElapsed(),
		)
	})

	// result
	result = &taskframework.TaskUnitRunResult{}
	muer.OnSuccess(func() {
		fmt.Printf("\n")
		fmt.Printf("[%s] 上传文件成功, 保存到网盘路径: %s\n", utu.taskInfo.Id(), utu.SavePath)
		// 统计
		utu.UploadStatistic.AddTotalSize(utu.LocalFileChecksum.Length)
		utu.UploadingDatabase.Delete(&utu.LocalFileChecksum.LocalFileMeta) // 删除
		utu.UploadingDatabase.Save()
		result.Succeed = true
	})
	muer.OnError(func(err error) {
		pcsError, ok := err.(pcserror.Error)
		if !ok {
			// 未知错误类型 (非预期的)
			// 不重试
			result.ResultMessage = "上传文件错误"
			result.Err = err
			return
		}

		// 默认需要重试
		result.NeedRetry = true

		switch pcsError.GetErrType() {
		case pcserror.ErrTypeRemoteError:
			// 远程百度服务器的错误
			switch pcsError.GetRemoteErrCode() {
			case 114514:
				// 自定义错误码, 仅在skip策略下出现
				result.ResultMessage = StrUploadFailed
				result.Err = pcsError
				if utu.Policy == baidupcs.SkipPolicy {
					result.Extra = baidupcs.SkipPolicy
					result.Err = nil
					result.ResultMessage = fmt.Sprintf("%s 目标已存在, 跳过", utu.SavePath)
				}
				result.NeedRetry = false
				return
			case 1919810:
				// 自定义错误码, 仅在rsync策略下出现
				result.Extra = baidupcs.RsyncPolicy
				result.Err = nil
				result.ResultMessage = fmt.Sprintf("%s 目标大小未发生改变, 跳过", utu.SavePath)
				result.NeedRetry = false
				return
			case 31363:
				// block miss in superfile2, 上传状态过期
				// 需要重试的
				utu.UploadingDatabase.Delete(&utu.LocalFileChecksum.LocalFileMeta)
				utu.UploadingDatabase.Save()

				result.ResultMessage = StrUploadFailed
				result.Err = errors.New("上传状态过期, 重新上传")
			case 31061:
				// 已存在重名文件, 不重试
				result.ResultMessage = StrUploadFailed
				result.Err = pcsError
				if utu.Policy == baidupcs.SkipPolicy {
					result.Extra = baidupcs.SkipPolicy
					result.Err = nil
					result.ResultMessage = fmt.Sprintf("%s 目标已存在, 跳过", utu.SavePath)
				}
				result.NeedRetry = false
				return
			default:
				result.ResultMessage = StrUploadFailed
				result.Err = pcsError
			}
		case pcserror.ErrTypeNetError:
			// 网络错误
			result.ResultMessage = StrUploadFailed
			result.Err = pcsError
			if strings.Contains(pcsError.GetError().Error(), "413 Request Entity Too Large") {
				// 请求实体过大
				// 不重试
				result.NeedRetry = false
				return
			}
		default:
			result.ResultMessage = StrUploadFailed
			result.NeedRetry = false
			result.Err = pcsError
		}
		return
	})
	muer.Execute()

	return
}

func (utu *UploadTaskUnit) OnRetry(lastRunResult *taskframework.TaskUnitRunResult) {
	// 输出错误信息
	if lastRunResult.Err == nil {
		// result中不包含Err, 忽略输出
		fmt.Printf("[%s] %s, 重试 %d/%d\n", utu.taskInfo.Id(), lastRunResult.ResultMessage, utu.taskInfo.Retry(), utu.taskInfo.MaxRetry())
		return
	}
	fmt.Printf("[%s] %s, %s, 重试 %d/%d\n", utu.taskInfo.Id(), lastRunResult.ResultMessage, lastRunResult.Err, utu.taskInfo.Retry(), utu.taskInfo.MaxRetry())
}

func (utu *UploadTaskUnit) OnSuccess(lastRunResult *taskframework.TaskUnitRunResult) {
}

func (utu *UploadTaskUnit) OnFailed(lastRunResult *taskframework.TaskUnitRunResult) {
	// 失败
	if lastRunResult.Err == nil {
		// result中不包含Err, 忽略输出
		fmt.Printf("[%s] %s\n", utu.taskInfo.Id(), lastRunResult.ResultMessage)
		return
	}
	fmt.Printf("[%s] %s, %s\n", utu.taskInfo.Id(), lastRunResult.ResultMessage, lastRunResult.Err)
}

func (utu *UploadTaskUnit) OnComplete(lastRunResult *taskframework.TaskUnitRunResult) {
}

func (utu *UploadTaskUnit) RetryWait() time.Duration {
	return pcsfunctions.RetryWait(utu.taskInfo.Retry())
}

func (utu *UploadTaskUnit) Run() (result *taskframework.TaskUnitRunResult) {
	fmt.Printf("[%s] 准备上传: %s\n", utu.taskInfo.Id(), utu.LocalFileChecksum.Path)

	if utu.LocalFileChecksum.Length > baidupcs.MaxUploadSize {
		fmt.Printf("[%s] 文件大小超过128G, 无法上传, 跳过...\n", utu.taskInfo.Id())
		return
	}

	err := utu.LocalFileChecksum.OpenPath()
	if err != nil {
		fmt.Printf("[%s] 文件不可读, 错误信息: %s, 跳过...\n", utu.taskInfo.Id(), err)
		return
	}
	defer utu.LocalFileChecksum.Close() // 关闭文件

	// 准备文件
	utu.prepareFile()

	switch utu.Step {
	case StepUploadRapidUpload:
		goto stepUploadRapidUpload
	case StepUploadUpload:
		goto stepUploadUpload
	case JustGoon:
		return
	}

stepUploadRapidUpload:
	// 秒传
	{
		isContinue, rapidUploadResult := utu.rapidUpload()
		if !isContinue {
			// 不继续, 返回秒传的结果
			return rapidUploadResult
		}
	}

stepUploadUpload:
	// 正常上传流程
	uploadResult := utu.upload()

	return uploadResult
}
