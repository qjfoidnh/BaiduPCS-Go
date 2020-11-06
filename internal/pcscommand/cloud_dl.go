package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
)

// RunCloudDlAddTask 执行添加离线下载任务
func RunCloudDlAddTask(sourceURLs []string, savePath string) {
	var (
		err error
		pcs = GetBaiduPCS()
	)
	err = matchPathByShellPatternOnce(&savePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var taskid int64
	for k := range sourceURLs {
		taskid, err = pcs.CloudDlAddTask(sourceURLs[k], savePath+baidupcs.PathSeparator)
		if err != nil {
			fmt.Printf("[%d] %s, 地址: %s\n", k+1, err, sourceURLs[k])
			continue
		}

		fmt.Printf("[%d] 添加离线任务成功, 任务ID(task_id): %d, 源地址: %s, 保存路径: %s\n", k+1, taskid, sourceURLs[k], savePath)
	}
}

// RunCloudDlQueryTask 精确查询离线下载任务
func RunCloudDlQueryTask(taskIDs []int64) {
	cl, err := GetBaiduPCS().CloudDlQueryTask(taskIDs)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Println(cl)
}

// RunCloudDlListTask 查询离线下载任务列表
func RunCloudDlListTask() {
	cl, err := GetBaiduPCS().CloudDlListTask()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Println(cl)
}

// RunCloudDlCancelTask 取消离线下载任务
func RunCloudDlCancelTask(taskIDs []int64) {
	for _, id := range taskIDs {
		err := GetBaiduPCS().CloudDlCancelTask(id)
		if err != nil {
			fmt.Printf("[%d] %s\n", id, err)
			continue
		}

		fmt.Printf("[%d] 取消成功\n", id)
	}
}

// RunCloudDlDeleteTask 删除离线下载任务
func RunCloudDlDeleteTask(taskIDs []int64) {
	for _, id := range taskIDs {
		err := GetBaiduPCS().CloudDlDeleteTask(id)
		if err != nil {
			fmt.Printf("[%d] %s\n", id, err)
			continue
		}

		fmt.Printf("[%d] 删除成功\n", id)
	}
}

// RunCloudDlClearTask 清空离线下载任务记录
func RunCloudDlClearTask() {
	total, err := GetBaiduPCS().CloudDlClearTask()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("%s成功, 共清除 %d 条记录\n", baidupcs.OperationCloudDlClearTask, total)
	return
}
