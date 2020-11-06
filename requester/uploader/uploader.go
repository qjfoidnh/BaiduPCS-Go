// Package uploader 上传包
package uploader

import (
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/BaiduPCS-Go/pcsverbose"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/BaiduPCS-Go/requester/rio"
	"net/http"
	"time"
)

const (
	// BufioReadSize bufio 缓冲区大小, 用于上传时读取文件
	BufioReadSize = int(64 * converter.KB) // 64KB
)

type (
	//CheckFunc 上传完成的检测函数
	CheckFunc func(resp *http.Response, uploadErr error)

	// Uploader 上传
	Uploader struct {
		url         string   // 上传地址
		readed64    Readed64 // 要上传的对象
		contentType string

		client *requester.HTTPClient

		executeTime time.Time
		executed    bool
		finished    chan struct{}

		checkFunc CheckFunc
		onExecute func()
		onFinish  func()
	}
)

var (
	uploaderVerbose = pcsverbose.New("UPLOADER")
)

// NewUploader 返回 uploader 对象, url: 上传地址, readerlen64: 实现 rio.ReaderLen64 接口的对象, 例如文件
func NewUploader(url string, readerlen64 rio.ReaderLen64) (uploader *Uploader) {
	uploader = &Uploader{
		url:      url,
		readed64: NewReaded64(readerlen64),
	}

	return
}

func (u *Uploader) lazyInit() {
	if u.finished == nil {
		u.finished = make(chan struct{})
	}
	if u.client == nil {
		u.client = requester.NewHTTPClient()
	}
	u.client.SetTimeout(0)
	u.client.SetResponseHeaderTimeout(0)
}

// SetClient 设置http客户端
func (u *Uploader) SetClient(c *requester.HTTPClient) {
	u.client = c
}

//SetContentType 设置Content-Type
func (u *Uploader) SetContentType(contentType string) {
	u.contentType = contentType
}

//SetCheckFunc 设置上传完成的检测函数
func (u *Uploader) SetCheckFunc(checkFunc CheckFunc) {
	u.checkFunc = checkFunc
}

// Execute 执行上传, 收到返回值信号则为上传结束
func (u *Uploader) Execute() {
	pcsutil.Trigger(u.onExecute)

	// 开始上传
	u.executeTime = time.Now()
	u.executed = true
	resp, _, err := u.execute()

	// 上传结束
	close(u.finished)

	if u.checkFunc != nil {
		u.checkFunc(resp, err)
	}

	pcsutil.Trigger(u.onFinish) // 触发上传结束的事件
}

func (u *Uploader) execute() (resp *http.Response, code int, err error) {
	u.lazyInit()
	header := map[string]string{}
	if u.contentType != "" {
		header["Content-Type"] = u.contentType
	}

	resp, err = u.client.Req(http.MethodPost, u.url, u.readed64, header)
	if err != nil {
		return nil, 2, err
	}

	return resp, 0, nil
}

// OnExecute 任务开始时触发的事件
func (u *Uploader) OnExecute(fn func()) {
	u.onExecute = fn
}

// OnFinish 任务完成时触发的事件
func (u *Uploader) OnFinish(fn func()) {
	u.onFinish = fn
}
