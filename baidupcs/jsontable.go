package baidupcs

import (
	"errors"
	"github.com/iikira/BaiduPCS-Go/pcstable"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/json-iterator/go"
	"path"
	"strconv"
	"strings"
)

type (
	// PathJSON 网盘路径
	PathJSON struct {
		Path string `json:"path"`
	}

	// PathsListJSON 网盘路径列表
	PathsListJSON struct {
		List []*PathJSON `json:"list"`
	}

	// FsIDJSON 文件或目录ID
	FsIDJSON struct {
		FsID int64 `json:"fs_id"` // fs_id
	}

	// FsIDListJSON fs_id 列表
	FsIDListJSON struct {
		List []*FsIDJSON `json:"list"`
	}

	// CpMvJSON 源文件目录的地址和目标文件目录的地址
	CpMvJSON struct {
		From string `json:"from"` // 源文件或目录
		To   string `json:"to"`   // 目标文件或目录
	}

	// CpMvJSONList CpMvJSON 列表
	CpMvJSONList []*CpMvJSON

	// CpMvListJSON []*CpMvJSON 对象数组
	CpMvListJSON struct {
		List CpMvJSONList `json:"list"`
	}

	// BlockListJSON 文件分块信息JSON
	BlockListJSON struct {
		BlockList []string `json:"block_list"`
	}
)

var (
	// ErrNilJSONValue 解析出的json值为空
	ErrNilJSONValue = errors.New("json value is nil")
)

// JSON json 数据构造
func (plj *PathsListJSON) JSON(paths ...string) (data []byte, err error) {
	plj.List = make([]*PathJSON, len(paths))

	for k := range paths {
		plj.List[k] = &PathJSON{
			Path: paths[k],
		}
	}

	data, err = jsoniter.Marshal(plj)
	return
}

// JSON json 数据构造
func (cj *CpMvJSON) JSON() (data []byte, err error) {
	data, err = jsoniter.Marshal(cj)
	return
}

// JSON json 数据构造
func (clj *CpMvListJSON) JSON() (data []byte, err error) {
	data, err = jsoniter.Marshal(clj)
	return
}

func (clj *CpMvListJSON) String() string {
	builder := &strings.Builder{}

	tb := pcstable.NewTable(builder)
	tb.SetHeader([]string{"#", "原路径", "目标路径"})

	for k := range clj.List {
		if clj.List[k] == nil {
			continue
		}
		tb.Append([]string{strconv.Itoa(k), clj.List[k].From, clj.List[k].To})
	}

	tb.Render()
	return builder.String()
}

// AllRelatedDir 获取所有相关的目录
func (cjl *CpMvJSONList) AllRelatedDir() (dirs []string) {
	for _, cj := range *cjl {
		fromDir, toDir := path.Dir(cj.From), path.Dir(cj.To)
		if !pcsutil.ContainsString(dirs, fromDir) {
			dirs = append(dirs, fromDir)
		}
		if !pcsutil.ContainsString(dirs, toDir) {
			dirs = append(dirs, toDir)
		}
	}
	return
}
