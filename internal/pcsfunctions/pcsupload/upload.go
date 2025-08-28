package pcsupload

import (
	"context"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
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

type PCSInfo struct {
	*pcserror.PCSErrInfo
	Host   string   `json:"host"`
	Server []string `json:"server"`
}

var client = pcsconfig.Config.PCSHTTPClient()

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

// Precreate 检查网盘的目标路径是否已存在同名文件及路径合法性, 顺便获取本次上传用的pcs服务器
func (pu *PCSUpload) Precreate(fileSize int64, policy string) (pcsHost string, pcsError pcserror.Error) {
	pcsError = pu.pcs.CheckIsdir(baidupcs.OperationUpload, pu.targetPath, policy, fileSize)
	if pcsError != nil {
		return
	}

	dataReadCloser, pcsError := pu.pcs.PreparePCSServers()
	if pcsError != nil {
		return
	}
	defer dataReadCloser.Close()
	pcsInfo := &PCSInfo{
		PCSErrInfo: pcserror.NewPCSErrorInfo(baidupcs.OperationGetPCSServer),
	}
	pcsError = pcserror.HandleJSONParse(baidupcs.OperationGetPCSServer, dataReadCloser, pcsInfo)
	if pcsError != nil {
		return
	}
	if len(pcsInfo.Server) > 0 {
		pcsHost = pcsInfo.Server[0]
	} else {
		pcsHost = pcsInfo.Host
	}

	return

}

func (pu *PCSUpload) TmpFile(ctx context.Context, pcsHost, uploadId, targetPath string, partSeq int, partOffset int64, r rio.ReaderLen64) (checksum string, uperr error) {
	pu.lazyInit()

	var respErr *uploader.MultiError

	// 临时切换为动态pcs addr
	originPCSHost := pu.pcs.GetPCSAddr()
	defer pu.pcs.SetPCSAddr(originPCSHost)
	pu.pcs.SetPCSAddr(pcsHost)

	checksum, pcsError := pu.pcs.UploadTmpFile(uploadId, targetPath, partSeq, partOffset, func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
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

func (pu *PCSUpload) CreateSuperFile(uploadId string, fileSize int64, checksumMap map[int]string) (err error) {
	pu.lazyInit()
	return pu.pcs.UploadCreateSuperFile(uploadId, fileSize, pu.targetPath, checksumMap)
}
