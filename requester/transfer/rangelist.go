package transfer

import (
	"errors"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"sync"
	"sync/atomic"
)

type (
	//RangeList 请求范围列表
	RangeList []*Range

	//RangeListGen Range 生成器
	RangeListGen struct {
		total        int64
		begin        int64
		blockSize    int64
		parallel     int
		count        int // 已生成次数
		rangeGenMode RangeGenMode
		mu           sync.Mutex
	}
)

const (
	// DefaultBlockSize 默认的BlockSize
	DefaultBlockSize = 256 * converter.KB
)

var (
	// ErrUnknownRangeGenMode RangeGenMode 非法
	ErrUnknownRangeGenMode = errors.New("Unknown RangeGenMode")
)

//Len 长度
func (r *Range) Len() int64 {
	return r.LoadEnd() - r.LoadBegin()
}

//LoadBegin 读取Begin, 原子操作
func (r *Range) LoadBegin() int64 {
	return atomic.LoadInt64(&r.Begin)
}

//AddBegin 增加Begin, 原子操作
func (r *Range) AddBegin(i int64) (newi int64) {
	return atomic.AddInt64(&r.Begin, i)
}

//LoadEnd 读取End, 原子操作
func (r *Range) LoadEnd() int64 {
	return atomic.LoadInt64(&r.End)
}

//StoreBegin 储存End, 原子操作
func (r *Range) StoreBegin(end int64) {
	atomic.StoreInt64(&r.Begin, end)
}

//StoreEnd 储存End, 原子操作
func (r *Range) StoreEnd(end int64) {
	atomic.StoreInt64(&r.End, end)
}

// ShowDetails 显示Range细节
func (r *Range) ShowDetails() string {
	return fmt.Sprintf("{%d-%d}", r.LoadBegin(), r.LoadEnd())
}

//Len 获取所有的Range的剩余长度
func (rl *RangeList) Len() int64 {
	var l int64
	for _, wrange := range *rl {
		if wrange == nil {
			continue
		}
		l += wrange.Len()
	}
	return l
}

// NewRangeListGenDefault 初始化默认Range生成器, 根据parallel平均生成
func NewRangeListGenDefault(totalSize, begin int64, count, parallel int) *RangeListGen {
	return &RangeListGen{
		total:        totalSize,
		begin:        begin,
		parallel:     parallel,
		count:        count,
		rangeGenMode: RangeGenMode_Default,
	}
}

// NewRangeListGenBlockSize 初始化Range生成器, 根据blockSize生成
func NewRangeListGenBlockSize(totalSize, begin, blockSize int64) *RangeListGen {
	return &RangeListGen{
		total:        totalSize,
		begin:        begin,
		blockSize:    blockSize,
		rangeGenMode: RangeGenMode_BlockSize,
	}
}

// RangeGenMode 返回Range生成方式
func (gen *RangeListGen) RangeGenMode() RangeGenMode {
	return gen.rangeGenMode
}

// RangeCount 返回预计生成的Range数量
func (gen *RangeListGen) RangeCount() (rangeCount int) {
	switch gen.rangeGenMode {
	case RangeGenMode_Default:
		rangeCount = gen.parallel - gen.count
		//rangeCount = int(math.Ceil(float64(gen.total) / float64(gen.blockSize)))
	case RangeGenMode_BlockSize:
		rangeCount = int((gen.total - gen.begin) / gen.blockSize)
		if gen.total%gen.blockSize != 0 {
			rangeCount++
		}
	}
	return
}

// LoadBegin 返回begin
func (gen *RangeListGen) LoadBegin() (begin int64) {
	gen.mu.Lock()
	begin = gen.begin
	gen.mu.Unlock()
	return
}

// LoadBlockSize 返回blockSize
func (gen *RangeListGen) LoadBlockSize() (blockSize int64) {
	switch gen.rangeGenMode {
	case RangeGenMode_Default:
		if gen.blockSize <= 0 {
			gen.blockSize = (gen.total - gen.begin) / int64(gen.parallel)
			if gen.blockSize < 256 * converter.KB {
				gen.blockSize = 256 * converter.KB
			}
		}
		blockSize = gen.blockSize
	case RangeGenMode_BlockSize:
		blockSize = gen.blockSize
	}
	return
}

// IsDone 是否已分配完成
func (gen *RangeListGen) IsDone() bool {
	return gen.begin >= gen.total
}

// GenRange 生成 Range
func (gen *RangeListGen) GenRange() (index int, r *Range) {
	var (
		end int64
	)
	if gen.parallel < 1 {
		gen.parallel = 1
	}
	switch gen.rangeGenMode {
	case RangeGenMode_Default:
		gen.LoadBlockSize()
		gen.mu.Lock()
		defer gen.mu.Unlock()

		if gen.IsDone() {
			return gen.count, nil
		}

		gen.count++
		if gen.count >= gen.parallel {
			end = gen.total
		} else {
			end = gen.begin + gen.blockSize
		}

		//end = gen.begin + gen.blockSize
		if end >= gen.total {
			end = gen.total
		}

		r = &Range{
			Begin: gen.begin,
			End:   end,
		}

		gen.begin = end
		index = gen.count - 1
		return
	case RangeGenMode_BlockSize:
		if gen.blockSize <= 0 {
			gen.blockSize = DefaultBlockSize
		}
		gen.mu.Lock()
		defer gen.mu.Unlock()

		if gen.IsDone() {
			return gen.count, nil
		}

		gen.count++
		end = gen.begin + gen.blockSize
		if end >= gen.total {
			end = gen.total
		}
		r = &Range{
			Begin: gen.begin,
			End:   end,
		}
		gen.begin = end
		index = gen.count - 1
		return
	}

	return 0, nil
}
