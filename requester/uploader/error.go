package uploader

type (
	// MultiError 多线程上传的错误
	MultiError struct {
		Err error
		// IsRetry 是否重试,
		Terminated bool
	}
)

func (me *MultiError) Error() string {
	return me.Err.Error()
}
