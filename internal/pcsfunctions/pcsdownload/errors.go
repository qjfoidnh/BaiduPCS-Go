package pcsdownload

import "errors"

var (
	// ErrDownloadNotSupportChecksum 文件不支持校验
	ErrDownloadNotSupportChecksum = errors.New("该文件不支持校验")
	// ErrDownloadChecksumFailed 文件校验失败
	ErrDownloadChecksumFailed = errors.New("该文件校验失败, 文件md5值与服务器记录的不匹配")
	// ErrDownloadFileBanned 违规文件
	ErrDownloadFileBanned = errors.New("该文件可能是违规文件, 不支持校验")
	// ErrDlinkNotFound 未取得下载链接
	ErrDlinkNotFound = errors.New("未取得下载链接")
	// ErrShareInfoNotFound 未在已分享列表中找到分享信息
	ErrShareInfoNotFound = errors.New("未在已分享列表中找到分享信息")
)
