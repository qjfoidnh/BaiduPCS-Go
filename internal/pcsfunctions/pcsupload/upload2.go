package pcsupload

/*
import (
	"github.com/iikira/BaiduPCS-Go/baidupcs"
)

type (
	// PCSUpload2 新的上传方式
	// TODO
	PCSUpload2 struct {
		pcs        *baidupcs.BaiduPCS
		targetPath string
		uploadid   string
	}
)

func NewPCSUpload2(pcs *baidupcs.BaiduPCS, targetPath string) uploader.MultiUpload {
	return &PCSUpload{
		pcs:        pcs,
		targetPath: targetPath,
	}
}

func (pu2 *PCSUpload2) lazyInit() {
	if pu2.pcs == nil {
		pu2.pcs = &baidupcs.BaiduPCS{}
	}
}

// Precreate
func (pu2 *PCSUpload2) Precreate() (err error) {
	return nil
}

func (pu2 *PCSUpload2) TmpFile(ctx context.Context, partseq int, partOffset int64, r rio.ReaderLen64) (checksum string, uperr error) {
	pu2.lazyInit()
	return pu.pcs.UploadTmpFile(func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
		client := pcsconfig.Config.HTTPClient()
		client.SetCookiejar(jar)
		client.SetTimeout(0)

		mr := multipartreader.NewMultipartReader()
		mr.AddFormFile("uploadedfile", "", r)
		mr.CloseMultipart()

		doneChan := make(chan struct{}, 1)
		go func() {
			resp, err = client.Req("POST", uploadURL, mr, nil)
			doneChan <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			return resp, ctx.Err()
		case <-doneChan:
			// return
		}
		return
	})
}

func (pu2 *PCSUpload2) CreateSuperFile(checksumList ...string) (err error) {
	return nil
}
*/
