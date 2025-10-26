package uploader

import (
	"context"
	"errors"
	"github.com/oleiade/lane"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/waitgroup"
	"os"
	"sync"
)

type (
	worker struct {
		id         int
		partOffset int64
		splitUnit  SplitUnit
		checksum   string
	}

	workerList []*worker
)

// CheckSumList 返回所以worker的checksum
// TODO: 实现sort
func (werl *workerList) CheckSumList() []string {
	checksumList := make([]string, 0, len(*werl))
	for _, wer := range *werl {
		checksumList = append(checksumList, wer.checksum)
	}
	return checksumList
}

func (werl *workerList) Readed() int64 {
	var readed int64
	for _, wer := range *werl {
		readed += wer.splitUnit.Readed()
	}
	return readed
}

func (muer *MultiUploader) upload() (uperr error) {
	originPCSHost, err := muer.multiUpload.Precreate()
	if err != nil {
		return err
	}
	var (
		uploadDeque = lane.NewDeque()
	)

	var (
		checksumMap = make(map[int]string) // key: wer.id, value: checksum
		mu          sync.Mutex
	)

	// 加入队列
	for _, wer := range muer.workers {
		if wer.checksum == "" {
			uploadDeque.Append(wer)
		}
	}

	for {
		wg := waitgroup.NewWaitGroup(muer.config.Parallel)
		for {
			e := uploadDeque.Shift()
			if e == nil { // 任务为空
				break
			}

			wer := e.(*worker)
			wg.AddDelta()
			go func() {
				defer wg.Done()

				var (
					ctx, cancel = context.WithCancel(context.Background())
					doneChan    = make(chan struct{})
					checksum    string
					terr        error
				)
				go func() {
					checksum, terr = muer.multiUpload.TmpFile(ctx, muer.instanceState.Uploadid, muer.targetPath, wer.id, wer.partOffset, wer.splitUnit)
					close(doneChan)
				}()
				select {
				case <-muer.canceled:
					cancel()
					return
				case <-doneChan:
					// continue
				}
				cancel()
				if terr != nil {
					var me *MultiError
					if errors.As(terr, &me) {
						if me.Terminated { // 终止
							muer.closeCanceledOnce.Do(func() { // 只关闭一次
								close(muer.canceled)
							})
							uperr = me.Err
							return
						}
					}

					uploaderVerbose.Warnf("upload err: %s, id: %d\n", terr, wer.id)
					wer.splitUnit.Seek(0, os.SEEK_SET)
					uploadDeque.Append(wer)
					return
				}
				wer.checksum = checksum
				mu.Lock()
				checksumMap[wer.id] = checksum // 记录成功任务的 checksum
				mu.Unlock()

				// 通知更新
				if muer.updateInstanceStateChan != nil && len(muer.updateInstanceStateChan) < cap(muer.updateInstanceStateChan) {
					muer.updateInstanceStateChan <- struct{}{}
				}
			}()
		}
		wg.Wait()

		// 没有任务了
		if uploadDeque.Size() == 0 {
			break
		}
	}

	select {
	case <-muer.canceled:
		if uperr != nil {
			return uperr
		}
		return context.Canceled
	default:
	}

	cerr := muer.multiUpload.CreateSuperFile(originPCSHost, muer.config.Policy, muer.instanceState.Uploadid, muer.file.Len(), checksumMap)
	if cerr != nil {
		return cerr
	}

	return
}
