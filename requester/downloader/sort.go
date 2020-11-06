package downloader

type (
	// ByLeftDesc 根据剩余下载量倒序排序
	ByLeftDesc struct {
		WorkerList
	}
)

// Len 返回长度
func (wl WorkerList) Len() int {
	return len(wl)
}

// Swap 交换
func (wl WorkerList) Swap(i, j int) {
	wl[i], wl[j] = wl[j], wl[i]
}

// Less 实现倒序
func (wl ByLeftDesc) Less(i, j int) bool {
	return wl.WorkerList[i].wrange.Len() > wl.WorkerList[j].wrange.Len()
}
