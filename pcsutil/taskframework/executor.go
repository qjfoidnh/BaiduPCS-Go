package taskframework

import (
	"github.com/GeertJohan/go.incremental"
	"github.com/oleiade/lane"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/waitgroup"
	"strconv"
	"time"
)

type (
	TaskExecutor struct {
		incr     *incremental.Int // 任务id生成
		deque    *lane.Deque      // 队列
		parallel int              // 任务的最大并发量

		// 是否统计失败队列
		IsFailedDeque bool
		failedDeque   *lane.Deque
	}
)

func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{}
}

func (te *TaskExecutor) lazyInit() {
	if te.deque == nil {
		te.deque = lane.NewDeque()
	}
	if te.incr == nil {
		te.incr = &incremental.Int{}
	}
	if te.parallel < 1 {
		te.parallel = 1
	}
	if te.IsFailedDeque {
		te.failedDeque = lane.NewDeque()
	}
}

// 设置任务的最大并发量
func (te *TaskExecutor) SetParallel(parallel int) {
	te.parallel = parallel
}

//Append 将任务加到任务队列末尾
func (te *TaskExecutor) Append(unit TaskUnit, maxRetry int) *TaskInfo {
	te.lazyInit()
	taskInfo := &TaskInfo{
		id:       strconv.Itoa(te.incr.Next()),
		maxRetry: maxRetry,
	}
	unit.SetTaskInfo(taskInfo)
	te.deque.Append(&TaskInfoItem{
		Info: taskInfo,
		Unit: unit,
	})
	return taskInfo
}

//AppendNoRetry 将任务加到任务队列末尾, 不重试
func (te *TaskExecutor) AppendNoRetry(unit TaskUnit) {
	te.Append(unit, 0)
}

//Count 返回任务数量
func (te *TaskExecutor) Count() int {
	if te.deque == nil {
		return 0
	}
	return te.deque.Size()
}

//Execute 执行任务
func (te *TaskExecutor) Execute() {
	te.lazyInit()

	for {
		wg := waitgroup.NewWaitGroup(te.parallel)
		for {
			e := te.deque.Shift()
			if e == nil { // 任务为空
				break
			}

			// 获取任务
			task := e.(*TaskInfoItem)
			wg.AddDelta()

			go func(task *TaskInfoItem) {
				defer wg.Done()

				result := task.Unit.Run()

				// 返回结果为空
				if result == nil {
					task.Unit.OnComplete(result)
					return
				}

				if result.Succeed {
					task.Unit.OnSuccess(result)
					task.Unit.OnComplete(result)
					return
				}

				// 需要进行重试
				if result.NeedRetry {
					// 重试次数超出限制
					// 执行失败
					if task.Info.IsExceedRetry() {
						task.Unit.OnFailed(result)
						if te.IsFailedDeque {
							// 加入失败队列
							te.failedDeque.Append(task)
						}
						task.Unit.OnComplete(result)
						return
					}
					task.Info.retry++         // 增加重试次数
					task.Unit.OnRetry(result) // 调用重试
					task.Unit.OnComplete(result)

					time.Sleep(task.Unit.RetryWait()) // 等待
					te.deque.Append(task)             // 重新加入队列末尾
					return
				}

				// 执行失败
				task.Unit.OnFailed(result)
				if te.IsFailedDeque && result.Extra != "skip" {
					// 加入失败队列
					te.failedDeque.Append(task)
				}
				task.Unit.OnComplete(result)
			}(task)
		}

		wg.Wait()

		// 没有任务了
		if te.deque.Size() == 0 {
			break
		}
	}
}

//FailedDeque 获取失败队列
func (te *TaskExecutor) FailedDeque() *lane.Deque {
	return te.failedDeque
}

//Stop 停止执行
func (te *TaskExecutor) Stop() {

}

//Pause 暂停执行
func (te *TaskExecutor) Pause() {

}

//Resume 恢复执行
func (te *TaskExecutor) Resume() {
}
