package uploader

import (
	"github.com/qjfoidnh/BaiduPCS-Go/requester/transfer"
)

type (
	// BlockState 文件区块信息
	BlockState struct {
		ID       int            `json:"id"`
		Range    transfer.Range `json:"range"`
		CheckSum string         `json:"checksum"`
	}

	// InstanceState 上传断点续传信息
	InstanceState struct {
		BlockList []*BlockState `json:"block_list"`
	}
)

func (muer *MultiUploader) getWorkerListByInstanceState(is *InstanceState) workerList {
	workers := make(workerList, 0, len(is.BlockList))
	for _, blockState := range is.BlockList {
		if blockState.CheckSum == "" {
			// 这里CheckSum与worker的checksum值相同但意义不同，可以代表一个块是否已经上传到了服务器
			slicemd5, length, err := muer.calcuSplitMD5(blockState.Range.Begin)
			if err != nil || blockState.Range.Len() != int64(length){
				slicemd5 = blockState.CheckSum
			}
			workers = append(workers, &worker{
				id:         blockState.ID,
				partOffset: blockState.Range.Begin,
				splitUnit:  NewBufioSplitUnit(muer.file, blockState.Range, muer.speedsStat, muer.rateLimit),
				//checksum:   blockState.CheckSum,
				checksum:   slicemd5,
			})
		} else {
			// 已经完成的, 也要加入 (可继续优化)
			workers = append(workers, &worker{
				id:         blockState.ID,
				partOffset: blockState.Range.Begin,
				splitUnit: &fileBlock{
					readRange: blockState.Range,
					readed:    blockState.Range.End - blockState.Range.Begin,
					readerAt:  muer.file,
				},
				checksum: blockState.CheckSum,
				uploaded: true,
			})
		}
	}
	return workers
}
