// Package requester 提供网络请求简便操作
package requester

const (
	// DefaultUserAgent 默认浏览器标识
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
)

var (
	// UserAgent 浏览器标识
	UserAgent = DefaultUserAgent
	// DefaultClient 默认 http 客户端
	DefaultClient = NewHTTPClient()
)

type (
	// ContentTyper Content-Type 接口
	ContentTyper interface {
		ContentType() string
	}

	// ContentLengther Content-Length 接口
	ContentLengther interface {
		ContentLength() int64
	}

	// Event 下载/上传任务运行时事件
	Event func()

	// EventOnError 任务出错运行时事件
	EventOnError func(err error)
)
