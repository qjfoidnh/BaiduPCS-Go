package taskframework

type (
	TaskInfo struct {
		id       string
		maxRetry int
		retry    int
	}

	TaskInfoItem struct {
		Info *TaskInfo
		Unit TaskUnit
	}
)

// IsExceedRetry 重试次数达到限制
func (t *TaskInfo) IsExceedRetry() bool {
	return t.retry >= t.maxRetry
}

func (t *TaskInfo) Id() string {
	return t.id
}

func (t *TaskInfo) MaxRetry() int {
	return t.maxRetry
}

func (t *TaskInfo) SetMaxRetry(maxRetry int) {
	t.maxRetry = maxRetry
}

func (t *TaskInfo) Retry() int {
	return t.retry
}
