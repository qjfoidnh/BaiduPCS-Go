package downloader

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsverbose"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/transfer"
	"github.com/json-iterator/go"
	"os"
	"sync"
)

type (
	//InstanceState 状态, 断点续传信息
	InstanceState struct {
		saveFile *os.File
		format   InstanceStateStorageFormat
		ii       transfer.DownloadInstanceInfoExporter
		mu       sync.Mutex
	}

	// InstanceStateStorageFormat 断点续传储存类型
	InstanceStateStorageFormat int
)

const (
	// InstanceStateStorageFormatJSON json 格式
	InstanceStateStorageFormatJSON = iota
	// InstanceStateStorageFormatProto3 protobuf 格式
	InstanceStateStorageFormatProto3
)

//NewInstanceState 初始化InstanceState
func NewInstanceState(saveFile *os.File, format InstanceStateStorageFormat) *InstanceState {
	return &InstanceState{
		saveFile: saveFile,
		format:   format,
	}
}

func (is *InstanceState) checkSaveFile() bool {
	return is.saveFile != nil
}

func (is *InstanceState) getSaveFileContents() []byte {
	if !is.checkSaveFile() {
		return nil
	}

	finfo, err := is.saveFile.Stat()
	if err != nil {
		panic(err)
	}

	size := finfo.Size()
	if size > 0xffffffff {
		panic("savePath too large")
	}
	intSize := int(size)

	buf := cachepool.RawMallocByteSlice(intSize)

	n, _ := is.saveFile.ReadAt(buf, 0)
	return buf[:n]
}

//Get 获取断点续传信息
func (is *InstanceState) Get() (eii *transfer.DownloadInstanceInfo) {
	if !is.checkSaveFile() {
		return nil
	}

	is.mu.Lock()
	defer is.mu.Unlock()

	contents := is.getSaveFileContents()
	if len(contents) <= 0 {
		return
	}

	is.ii = &transfer.DownloadInstanceInfoExport{}
	var err error
	switch is.format {
	case InstanceStateStorageFormatProto3:
		err = proto.Unmarshal(contents, is.ii.(*transfer.DownloadInstanceInfoExport))
	default:
		err = jsoniter.Unmarshal(contents, is.ii)
	}

	if err != nil {
		pcsverbose.Verbosef("DEBUG: InstanceInfo unmarshal error: %s\n", err)
		return
	}

	eii = is.ii.GetInstanceInfo()
	return
}

//Put 提交断点续传信息
func (is *InstanceState) Put(eii *transfer.DownloadInstanceInfo) {
	if !is.checkSaveFile() {
		return
	}

	is.mu.Lock()
	defer is.mu.Unlock()

	if is.ii == nil {
		is.ii = &transfer.DownloadInstanceInfoExport{}
	}
	is.ii.SetInstanceInfo(eii)
	var (
		data []byte
		err  error
	)
	switch is.format {
	case InstanceStateStorageFormatProto3:
		data, err = proto.Marshal(is.ii.(*transfer.DownloadInstanceInfoExport))
	default:
		data, err = jsoniter.Marshal(is.ii)
	}
	if err != nil {
		panic(err)
	}

	err = is.saveFile.Truncate(int64(len(data)))
	if err != nil {
		pcsverbose.Verbosef("DEBUG: truncate file error: %s\n", err)
	}

	_, err = is.saveFile.WriteAt(data, 0)
	if err != nil {
		pcsverbose.Verbosef("DEBUG: write instance state error: %s\n", err)
	}
}

//Close 关闭
func (is *InstanceState) Close() error {
	if !is.checkSaveFile() {
		return nil
	}

	return is.saveFile.Close()
}

func (der *Downloader) initInstanceState(format InstanceStateStorageFormat) (err error) {
	if der.instanceState != nil {
		return errors.New("already initInstanceState")
	}

	var saveFile *os.File
	if !der.config.IsTest && der.config.InstanceStatePath != "" {
		saveFile, err = os.OpenFile(der.config.InstanceStatePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
	}

	der.instanceState = NewInstanceState(saveFile, format)
	return nil
}

func (der *Downloader) removeInstanceState() error {
	der.instanceState.Close()
	if !der.config.IsTest && der.config.InstanceStatePath != "" {
		return os.Remove(der.config.InstanceStatePath)
	}
	return nil
}
