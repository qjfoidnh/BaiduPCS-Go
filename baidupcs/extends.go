package baidupcs

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/iikira/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/iikira/BaiduPCS-Go/pcsutil/escaper"
	"github.com/iikira/BaiduPCS-Go/requester/downloader"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const (
	// ShellPatternCharacters 通配符字符串
	ShellPatternCharacters = "*?[]"
)

var (
	// ErrFixMD5Isdir 目录不需要修复md5
	ErrFixMD5Isdir = errors.New("directory not support fix md5")
	// ErrFixMD5Failed 修复MD5失败, 可能服务器未刷新
	ErrFixMD5Failed = errors.New("fix md5 failed")
	// ErrFixMD5FileInfoNil 文件信息对象为空
	ErrFixMD5FileInfoNil = errors.New("file info is nil")
	// ErrMatchPathByShellPatternNotAbsPath 不是绝对路径
	ErrMatchPathByShellPatternNotAbsPath = errors.New("not absolute path")

	ErrContentRangeNotFound               = errors.New("Content-Range not found")
	ErrGetRapidUploadInfoLengthNotFound   = errors.New("Content-Length not found")
	ErrGetRapidUploadInfoMD5NotFound      = errors.New("Content-MD5 not found")
	ErrGetRapidUploadInfoCrc32NotFound    = errors.New("x-bs-meta-crc32 not found")
	ErrGetRapidUploadInfoFilenameNotEqual = errors.New("文件名不匹配")
	ErrGetRapidUploadInfoLengthNotEqual   = errors.New("Content-Length 不匹配")
	ErrGetRapidUploadInfoMD5NotEqual      = errors.New("Content-MD5 不匹配")
	ErrGetRapidUploadInfoCrc32NotEqual    = errors.New("x-bs-meta-crc32 不匹配")
	ErrGetRapidUploadInfoSliceMD5NotEqual = errors.New("slice-md5 不匹配")

	ErrFileTooLarge = errors.New("文件大于20GB, 无法秒传")
)

func (pcs *BaiduPCS) getLocateDownloadLink(pcspath string) (link string, pcsError pcserror.Error) {
	info, pcsError := pcs.LocateDownload(pcspath)
	if pcsError != nil {
		return
	}

	u := info.SingleURL(pcs.isHTTPS)
	if u == nil {
		return "", &pcserror.PCSErrInfo{
			Operation: OperationLocateDownload,
			ErrType:   pcserror.ErrTypeOthers,
			Err:       ErrLocateDownloadURLNotFound,
		}
	}
	return u.String(), nil
}

// ExportByFileInfo 通过文件信息对象, 导出文件信息
func (pcs *BaiduPCS) ExportByFileInfo(finfo *FileDirectory) (rinfo *RapidUploadInfo, pcsError pcserror.Error) {
	errInfo := pcserror.NewPCSErrorInfo(OperationExportFileInfo)
	errInfo.ErrType = pcserror.ErrTypeOthers
	if finfo.Size > MaxRapidUploadSize {
		errInfo.Err = ErrFileTooLarge
		return nil, errInfo
	}

	rinfo, pcsError = pcs.GetRapidUploadInfoByFileInfo(finfo)
	if pcsError != nil {
		return nil, pcsError
	}
	if rinfo.Filename != finfo.Filename {
		baiduPCSVerbose.Infof("%s filename not equal, local: %s, remote link: %s\n", OperationExportFileInfo, finfo.Filename, rinfo.Filename)
		rinfo.Filename = finfo.Filename
	}
	return rinfo, nil
}

// GetRapidUploadInfoByFileInfo 通过文件信息对象, 获取秒传信息
func (pcs *BaiduPCS) GetRapidUploadInfoByFileInfo(finfo *FileDirectory) (rinfo *RapidUploadInfo, pcsError pcserror.Error) {
	if finfo.Size <= SliceMD5Size && len(finfo.BlockList) == 1 && finfo.BlockList[0] == finfo.MD5 {
		// 可直接秒传
		return &RapidUploadInfo{
			Filename:      finfo.Filename,
			ContentLength: finfo.Size,
			ContentMD5:    finfo.MD5,
			SliceMD5:      finfo.MD5,
			ContentCrc32:  "0",
		}, nil
	}

	link, pcsError := pcs.getLocateDownloadLink(finfo.Path)
	if pcsError != nil {
		return nil, pcsError
	}

	// 只有ContentLength可以比较
	// finfo记录的ContentMD5不一定是正确的
	// finfo记录的Filename不一定与获取到的一致
	return pcs.GetRapidUploadInfoByLink(link, &RapidUploadInfo{
		ContentLength: finfo.Size,
	})
}

// GetRapidUploadInfoByLink 通过下载链接, 获取文件秒传信息
func (pcs *BaiduPCS) GetRapidUploadInfoByLink(link string, compareRInfo *RapidUploadInfo) (rinfo *RapidUploadInfo, pcsError pcserror.Error) {
	errInfo := pcserror.NewPCSErrorInfo(OperationGetRapidUploadInfo)
	errInfo.ErrType = pcserror.ErrTypeOthers

	var (
		header     = pcs.getPanUAHeader()
		isSetRange = compareRInfo != nil && compareRInfo.ContentLength > SliceMD5Size // 是否设置Range
	)
	if isSetRange {
		header["Range"] = "bytes=0-" + strconv.FormatInt(SliceMD5Size-1, 10)
	}

	resp, err := pcs.client.Req(http.MethodGet, link, nil, header)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		errInfo.SetNetError(err)
		return nil, errInfo
	}

	// 检测响应状态码
	if resp.StatusCode/100 != 2 {
		errInfo.SetNetError(errors.New(resp.Status))
		return nil, errInfo
	}

	// 检测是否存在MD5
	md5Str := resp.Header.Get("Content-MD5")
	if md5Str == "" { // 未找到md5值, 可能是服务器未刷新
		errInfo.Err = ErrGetRapidUploadInfoMD5NotFound
		return nil, errInfo
	}
	if compareRInfo != nil && compareRInfo.ContentMD5 != "" && compareRInfo.ContentMD5 != md5Str {
		errInfo.Err = ErrGetRapidUploadInfoMD5NotEqual
		return nil, errInfo
	}

	// 获取文件名
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		errInfo.Err = err
		return nil, errInfo
	}
	filename, err := url.QueryUnescape(params["filename"])
	if err != nil {
		errInfo.Err = err
		return nil, errInfo
	}
	if compareRInfo != nil && compareRInfo.Filename != "" && compareRInfo.Filename != filename {
		errInfo.Err = ErrGetRapidUploadInfoFilenameNotEqual
		return nil, errInfo
	}

	var (
		contentLength int64
	)
	if isSetRange {
		// 检测Content-Range
		contentRange := resp.Header.Get("Content-Range")
		if contentRange == "" {
			errInfo.Err = ErrContentRangeNotFound
			return nil, errInfo
		}
		contentLength = downloader.ParseContentRange(contentRange)
	} else {
		contentLength = resp.ContentLength
	}

	// 检测Content-Length
	switch contentLength {
	case -1:
		errInfo.Err = ErrGetRapidUploadInfoLengthNotFound
		return nil, errInfo
	case 0:
		return &RapidUploadInfo{
			Filename:      filename,
			ContentLength: contentLength,
			ContentMD5:    EmptyContentMD5,
			SliceMD5:      EmptyContentMD5,
			ContentCrc32:  "0",
		}, nil
	default:
		if compareRInfo != nil && compareRInfo.ContentLength > 0 && compareRInfo.ContentLength != contentLength {
			errInfo.Err = ErrGetRapidUploadInfoLengthNotEqual
			return nil, errInfo
		}
	}

	// 检测是否存在crc32 值, 一般都会存在的
	crc32Str := resp.Header.Get("x-bs-meta-crc32")
	if crc32Str == "" || crc32Str == "0" {
		errInfo.Err = ErrGetRapidUploadInfoCrc32NotFound
		return nil, errInfo
	}
	if compareRInfo != nil && compareRInfo.ContentCrc32 != "" && compareRInfo.ContentCrc32 != crc32Str {
		errInfo.Err = ErrGetRapidUploadInfoCrc32NotEqual
		return nil, errInfo
	}

	// 获取slice-md5
	// 忽略比较slice-md5
	if contentLength <= SliceMD5Size {
		return &RapidUploadInfo{
			Filename:      filename,
			ContentLength: contentLength,
			ContentMD5:    md5Str,
			SliceMD5:      md5Str,
			ContentCrc32:  crc32Str,
		}, nil
	}

	buf := cachepool.RawMallocByteSlice(int(SliceMD5Size))
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		errInfo.SetNetError(err)
		return nil, errInfo
	}

	// 计算slice-md5
	m := md5.New()
	_, err = m.Write(buf)
	if err != nil {
		panic(err)
	}

	sliceMD5Str := hex.EncodeToString(m.Sum(nil))

	// 检测slice-md5, 不必要的
	if compareRInfo != nil && compareRInfo.SliceMD5 != "" && compareRInfo.SliceMD5 != sliceMD5Str {
		errInfo.Err = ErrGetRapidUploadInfoSliceMD5NotEqual
		return nil, errInfo
	}

	return &RapidUploadInfo{
		Filename:      filename,
		ContentLength: contentLength,
		ContentMD5:    md5Str,
		SliceMD5:      sliceMD5Str,
		ContentCrc32:  crc32Str,
	}, nil
}

// FixMD5ByFileInfo 尝试修复文件的md5, 通过文件信息对象
func (pcs *BaiduPCS) FixMD5ByFileInfo(finfo *FileDirectory) (pcsError pcserror.Error) {
	errInfo := pcserror.NewPCSErrorInfo(OperationFixMD5)
	errInfo.ErrType = pcserror.ErrTypeOthers
	if finfo == nil {
		errInfo.Err = ErrFixMD5FileInfoNil
		return errInfo
	}

	if finfo.Size > MaxRapidUploadSize { // 文件大于20GB
		errInfo.Err = ErrFileTooLarge
		return errInfo
	}

	// 忽略目录
	if finfo.Isdir {
		errInfo.Err = ErrFixMD5Isdir
		return errInfo
	}

	if len(finfo.BlockList) == 1 && finfo.BlockList[0] == finfo.MD5 {
		// 不需要修复
		return nil
	}

	link, pcsError := pcs.getLocateDownloadLink(finfo.Path)
	if pcsError != nil {
		return pcsError
	}

	var (
		cmpInfo = &RapidUploadInfo{
			Filename:      finfo.Filename,
			ContentLength: finfo.Size,
		}
	)
	rinfo, pcsError := pcs.GetRapidUploadInfoByLink(link, cmpInfo)
	if pcsError != nil {
		switch pcsError.GetError() {
		case ErrGetRapidUploadInfoMD5NotFound, ErrGetRapidUploadInfoCrc32NotFound:
			errInfo.Err = ErrFixMD5Failed
		default:
			errInfo.Err = pcsError
		}
		return errInfo
	}

	// 开始修复
	return pcs.RapidUploadNoCheckDir(finfo.Path, rinfo.ContentMD5, rinfo.SliceMD5, rinfo.ContentCrc32, rinfo.ContentLength)
}

// FixMD5 尝试修复文件的md5
func (pcs *BaiduPCS) FixMD5(pcspath string) (pcsError pcserror.Error) {
	finfo, pcsError := pcs.FilesDirectoriesMeta(pcspath)
	if pcsError != nil {
		return
	}

	return pcs.FixMD5ByFileInfo(finfo)
}

func (pcs *BaiduPCS) recurseMatchPathByShellPattern(index int, patternSlice *[]string, ps *[]string, pcspaths *[]string) {
	if index == len(*patternSlice) {
		*pcspaths = append(*pcspaths, strings.Join(*ps, PathSeparator))
		return
	}

	if !strings.ContainsAny((*patternSlice)[index], ShellPatternCharacters) {
		(*ps)[index] = (*patternSlice)[index]
		pcs.recurseMatchPathByShellPattern(index+1, patternSlice, ps, pcspaths)
		return
	}

	fds, pcsError := pcs.FilesDirectoriesList(strings.Join((*ps)[:index], PathSeparator), DefaultOrderOptions)
	if pcsError != nil {
		panic(pcsError) // 抛出异常
	}

	for k := range fds {
		if matched, _ := path.Match((*patternSlice)[index], fds[k].Filename); matched {
			(*ps)[index] = fds[k].Filename
			pcs.recurseMatchPathByShellPattern(index+1, patternSlice, ps, pcspaths)
		}
	}
	return
}

// MatchPathByShellPattern 通配符匹配文件路径, pattern 为绝对路径
func (pcs *BaiduPCS) MatchPathByShellPattern(pattern string) (pcspaths []string, pcsError pcserror.Error) {
	errInfo := pcserror.NewPCSErrorInfo(OperrationMatchPathByShellPattern)
	errInfo.ErrType = pcserror.ErrTypeOthers

	patternSlice := strings.Split(escaper.Escape(path.Clean(pattern), []rune{'['}), PathSeparator) // 转义中括号
	if patternSlice[0] != "" {
		errInfo.Err = ErrMatchPathByShellPatternNotAbsPath
		return nil, errInfo
	}

	ps := make([]string, len(patternSlice))
	defer func() { // 捕获异常
		if err := recover(); err != nil {
			pcspaths = nil
			pcsError = err.(pcserror.Error)
		}
	}()
	pcs.recurseMatchPathByShellPattern(1, &patternSlice, &ps, &pcspaths)
	return pcspaths, nil
}
