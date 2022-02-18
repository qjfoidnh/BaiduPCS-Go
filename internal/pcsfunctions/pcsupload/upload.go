package pcsupload

import (
	"context"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/multipartreader"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/rio"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/uploader"
	"io"
	"net/http"
	"time"
)

type (
	PCSUpload struct {
		pcs        *baidupcs.BaiduPCS
		targetPath string
	}

	EmptyReaderLen64 struct {
	}
)

var client *requester.HTTPClient = pcsconfig.Config.PCSHTTPClient()

func (e EmptyReaderLen64) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (e EmptyReaderLen64) Len() int64 {
	return 0
}

func NewPCSUpload(pcs *baidupcs.BaiduPCS, targetPath string) uploader.MultiUpload {
	return &PCSUpload{
		pcs:        pcs,
		targetPath: targetPath,
	}
}

func (pu *PCSUpload) lazyInit() {
	if pu.pcs == nil {
		pu.pcs = &baidupcs.BaiduPCS{}
	}
}

// Precreate 检查网盘的目标路径是否已存在同名文件及路径合法性
func (pu *PCSUpload) Precreate(fileSize int64, policy string) pcserror.Error {
	pcsError := pu.pcs.CheckIsdir(baidupcs.OperationUpload, pu.targetPath, policy, fileSize)
	return pcsError
}

func (pu *PCSUpload) TmpFile(ctx context.Context, partseq int, partOffset int64, r rio.ReaderLen64) (checksum string, uperr error) {
	pu.lazyInit()

	var respErr *uploader.MultiError
	checksum, pcsError := pu.pcs.UploadTmpFile(func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
		//client := pcsconfig.Config.PCSHTTPClient()
		client.SetCookiejar(jar)
		client.SetTimeout(200 * time.Second)

		mr := multipartreader.NewMultipartReader()
		mr.AddFormFile("uploadedfile", "", r)
		mr.CloseMultipart()

		doneChan := make(chan struct{}, 1)
		go func() {
			resp, err = client.Req(http.MethodPost, uploadURL, mr, nil)
			doneChan <- struct{}{}

			if resp != nil {
				// 不可恢复的错误
				switch resp.StatusCode {
				case 400, 401, 403, 413: // 4xx通常是由客户端非法操作引发，直接深度重试
					respErr = &uploader.MultiError{
						Terminated: true,
					}
				}
			}
		}()
		select {
		case <-ctx.Done(): // 取消
			// 返回, 让那边关闭连接
			return resp, ctx.Err()
		case <-doneChan:
			// return
		}
		return
	})

	if respErr != nil {
		respErr.Err = pcsError
		return checksum, respErr
	}

	return checksum, pcsError
}

func (pu *PCSUpload) CreateSuperFile(policy string, checksumList ...string) (err error) {
	pu.lazyInit()
	//newpath := ""
	// 先在网盘目标位置, 上传一个空文件
	// 防止出现file does not exist
	pcsError, newpath := pu.pcs.Upload(policy, pu.targetPath, func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
		mr := multipartreader.NewMultipartReader()
		mr.AddFormFile("file", "file", &EmptyReaderLen64{})
		mr.CloseMultipart()

		c := requester.NewHTTPClient()
		c.SetCookiejar(jar)
		return c.Req(http.MethodPost, uploadURL, mr, nil)
	})
	if pcsError != nil {
		// 修改操作
		pcsError.(*pcserror.PCSErrInfo).Operation = baidupcs.OperationUploadCreateSuperFile
		return pcsError
	}

	// 此时已到了最后的合并环节，policy只能使用overwrite, newpath而不用pu.targetPath是因为newcopy策略可能导致文件名变化
	//return pu.pcs.UploadCreateSuperFile("overwrite",false, pu.targetPath, checksumList...)
	return pu.pcs.UploadCreateSuperFile("overwrite",false, newpath, checksumList...)
}
