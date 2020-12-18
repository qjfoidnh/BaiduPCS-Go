package taskframework_test

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/taskframework"
	"testing"
	"time"
)

type (
	TestUnit struct {
		retry    bool
		taskInfo *taskframework.TaskInfo
	}
)

func (tu *TestUnit) SetTaskInfo(taskInfo *taskframework.TaskInfo) {
	tu.taskInfo = taskInfo
}

func (tu *TestUnit) OnFailed(lastRunResult *taskframework.TaskUnitRunResult) {
	fmt.Printf("[%s] error: %s, failed\n", tu.taskInfo.Id(), lastRunResult.Err)
}

func (tu *TestUnit) OnSuccess(lastRunResult *taskframework.TaskUnitRunResult) {
	fmt.Printf("[%s] success\n", tu.taskInfo.Id())
}

func (tu *TestUnit) OnComplete(lastRunResult *taskframework.TaskUnitRunResult) {
	fmt.Printf("[%s] complete\n", tu.taskInfo.Id())
}

func (tu *TestUnit) Run() (result *taskframework.TaskUnitRunResult) {
	fmt.Printf("[%s] running...\n", tu.taskInfo.Id())
	return &taskframework.TaskUnitRunResult{
		//Succeed:   true,
		NeedRetry: true,
	}
}

func (tu *TestUnit) OnRetry(lastRunResult *taskframework.TaskUnitRunResult) {
	fmt.Printf("[%s] prepare retry, times [%d/%d]...\n", tu.taskInfo.Id(), tu.taskInfo.Retry(), tu.taskInfo.MaxRetry())
}

func (tu *TestUnit) RetryWait() time.Duration {
	return 1 * time.Second
}

func TestTaskExecutor(t *testing.T) {
	te := taskframework.NewTaskExecutor()
	te.SetParallel(2)
	for i := 0; i < 3; i++ {
		tu := TestUnit{
			retry: false,
		}
		te.Append(&tu, 2)
	}
	te.Execute()
}
