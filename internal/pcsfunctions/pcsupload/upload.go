package pcsupload

import (
	"context"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/BaiduPCS-Go/requester/multipartreader"
	"github.com/iikira/BaiduPCS-Go/requester/rio"
	"github.com/iikira/BaiduPCS-Go/requester/uploader"
	"io"
	"net/http"
)

type (
	PCSUpload struct {
		pcs        *baidupcs.BaiduPCS
		targetPath string
	}

	EmptyReaderLen64 struct {
	}
)

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

// Precreate do nothing
func (pu *PCSUpload) Precreate() (err error) {
	return nil
}

func (pu *PCSUpload) TmpFile(ctx context.Context, partseq int, partOffset int64, r rio.ReaderLen64) (checksum string, uperr error) {
	pu.lazyInit()

	var respErr *uploader.MultiError
	checksum, pcsError := pu.pcs.UploadTmpFile(func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
		client := pcsconfig.Config.PCSHTTPClient()
		client.SetCookiejar(jar)
		client.SetTimeout(0)

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
				case 400, 401, 403, 413:
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

func (pu *PCSUpload) CreateSuperFile(checksumList ...string) (err error) {
	pu.lazyInit()

	// 先在网盘目标位置, 上传一个空文件
	// 防止出现file does not exist
	pcsError := pu.pcs.Upload(pu.targetPath, func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
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

	return pu.pcs.UploadCreateSuperFile(false, pu.targetPath, checksumList...)
}
