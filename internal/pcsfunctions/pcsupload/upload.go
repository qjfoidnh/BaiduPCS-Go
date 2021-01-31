package pcsupload

import (
	"context"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/multipartreader"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/rio"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/uploader"
	"io"
	"net/http"
)

type (
	PCSUpload struct {
		pcs        *baidupcs.BaiduPCS
		targetPath string
		localfilechecksum *checksum.LocalFileChecksum
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

func NewPCSUpload(pcs *baidupcs.BaiduPCS, targetPath string, filechecksum *checksum.LocalFileChecksum) uploader.MultiUpload {
	return &PCSUpload{
		pcs:        pcs,
		targetPath: targetPath,
		localfilechecksum: filechecksum,
	}
}

func (pu *PCSUpload) lazyInit() {
	if pu.pcs == nil {
		pu.pcs = &baidupcs.BaiduPCS{}
	}
}

// Precreate
func (pu *PCSUpload) Precreate(checksumList ...string) (precreateInfo *baidupcs.PrecreateInfo, err error) {
	//pu.pcs.UploadPrecreate(pu.targetPath, hex.EncodeToString(pu.localfilechecksum.MD5), hex.EncodeToString(pu.localfilechecksum.SliceMD5), strconv.FormatInt(int64(pu.localfilechecksum.CRC32), 10), pu.localfilechecksum.Length, checksumList...)
	precreateInfo, err = pu.pcs.UploadPrecreate(pu.targetPath, checksumList...)
	return

}

func (pu *PCSUpload) TmpFile(ctx context.Context, partseq int, partOffset int64, r rio.ReaderLen64, uploadid string) (checksum string, uperr error) {
	pu.lazyInit()

	var respErr *uploader.MultiError
	checksum, pcsError := pu.pcs.UploadSuperfile2(uploadid, pu.targetPath, partseq, partOffset, func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
		client := pcsconfig.Config.PCSHTTPClient()
		client.SetCookiejar(jar)
		client.SetTimeout(0)

		mr := multipartreader.NewMultipartReader()
		mr.AddFormFile("file", "", r)
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
	//checksum, pcsError := pu.pcs.UploadTmpFile(func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
	//	client := pcsconfig.Config.PCSHTTPClient()
	//	client.SetCookiejar(jar)
	//	client.SetTimeout(0)
	//
	//	mr := multipartreader.NewMultipartReader()
	//	mr.AddFormFile("uploadedfile", "", r)
	//	mr.CloseMultipart()
	//
	//	doneChan := make(chan struct{}, 1)
	//	go func() {
	//		resp, err = client.Req(http.MethodPost, uploadURL, mr, nil)
	//		doneChan <- struct{}{}
	//
	//		if resp != nil {
	//			// 不可恢复的错误
	//			switch resp.StatusCode {
	//			case 400, 401, 403, 413:
	//				respErr = &uploader.MultiError{
	//					Terminated: true,
	//				}
	//			}
	//		}
	//	}()
	//	select {
	//	case <-ctx.Done(): // 取消
	//		// 返回, 让那边关闭连接
	//		return resp, ctx.Err()
	//	case <-doneChan:
	//		// return
	//	}
	//	return
	//})

	if respErr != nil {
		respErr.Err = pcsError
		return checksum, respErr
	}

	return checksum, pcsError
}

func (pu *PCSUpload) CreateSuperFile(uploadid, bdstoken string, size int64, checksumList ...string) (err error) {
	pu.lazyInit()
	// 先在网盘目标位置, 上传一个空文件
	// 防止出现file does not exist
	//pcsError := pu.pcs.Upload(pu.targetPath, func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
	//	mr := multipartreader.NewMultipartReader()
	//	mr.AddFormFile("file", "file", &EmptyReaderLen64{})
	//	mr.CloseMultipart()
	//
	//	c := requester.NewHTTPClient()
	//	c.SetCookiejar(jar)
	//	return c.Req(http.MethodPost, uploadURL, mr, nil)
	//})
	//if pcsError != nil {
	//	// 修改操作
	//	pcsError.(*pcserror.PCSErrInfo).Operation = baidupcs.OperationUploadCreateSuperFile
	//	return pcsError
	//}
	//return pu.pcs.UploadCreateSuperFile(false, pu.targetPath, checksumList...)
	return pu.pcs.UploadCreateFile(uploadid, bdstoken, pu.targetPath, size, checksumList...)
}
