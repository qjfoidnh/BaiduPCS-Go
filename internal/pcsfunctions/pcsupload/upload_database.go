package pcsupload

import (
	"errors"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
	"github.com/iikira/BaiduPCS-Go/pcsutil/checksum"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/BaiduPCS-Go/pcsutil/jsonhelper"
	"github.com/iikira/BaiduPCS-Go/requester/uploader"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	// Uploading 未完成上传的信息
	Uploading struct {
		*checksum.LocalFileMeta
		State *uploader.InstanceState `json:"state"`
	}

	// UploadingDatabase 未完成上传的数据库
	UploadingDatabase struct {
		UploadingList []*Uploading `json:"upload_state"`
		Timestamp     int64        `json:"timestamp"`

		dataFile *os.File
	}
)

// NewUploadingDatabase 初始化未完成上传的数据库, 从库中读取内容
func NewUploadingDatabase() (ud *UploadingDatabase, err error) {
	file, err := os.OpenFile(filepath.Join(pcsconfig.GetConfigDir(), UploadingFileName), os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	ud = &UploadingDatabase{
		dataFile: file,
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() <= 0 {
		return ud, nil
	}

	err = jsonhelper.UnmarshalData(file, ud)
	if err != nil {
		return nil, err
	}

	return ud, nil
}

// Save 保存内容
func (ud *UploadingDatabase) Save() error {
	if ud.dataFile == nil {
		return errors.New("dataFile is nil")
	}

	ud.Timestamp = time.Now().Unix()

	var (
		builder = &strings.Builder{}
		err     = jsonhelper.MarshalData(builder, ud)
	)
	if err != nil {
		panic(err)
	}

	err = ud.dataFile.Truncate(int64(builder.Len()))
	if err != nil {
		return err
	}

	str := builder.String()
	_, err = ud.dataFile.WriteAt(converter.ToBytes(str), 0)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUploading 更新正在上传
func (ud *UploadingDatabase) UpdateUploading(meta *checksum.LocalFileMeta, state *uploader.InstanceState) {
	if meta == nil {
		return
	}

	meta.CompleteAbsPath()
	for k, uploading := range ud.UploadingList {
		if uploading.LocalFileMeta == nil {
			continue
		}
		if uploading.LocalFileMeta.EqualLengthMD5(meta) || uploading.LocalFileMeta.Path == meta.Path {
			ud.UploadingList[k].State = state
			return
		}
	}

	ud.UploadingList = append(ud.UploadingList, &Uploading{
		LocalFileMeta: meta,
		State:         state,
	})
}

func (ud *UploadingDatabase) deleteIndex(k int) {
	ud.UploadingList = append(ud.UploadingList[:k], ud.UploadingList[k+1:]...)
}

// Delete 删除
func (ud *UploadingDatabase) Delete(meta *checksum.LocalFileMeta) bool {
	if meta == nil {
		return false
	}

	meta.CompleteAbsPath()
	for k, uploading := range ud.UploadingList {
		if uploading.LocalFileMeta == nil {
			continue
		}
		if uploading.LocalFileMeta.EqualLengthMD5(meta) || uploading.LocalFileMeta.Path == meta.Path {
			ud.deleteIndex(k)
			return true
		}
	}
	return false
}

// Search 搜索
func (ud *UploadingDatabase) Search(meta *checksum.LocalFileMeta) *uploader.InstanceState {
	if meta == nil {
		return nil
	}

	meta.CompleteAbsPath()
	ud.clearModTimeChange()
	for _, uploading := range ud.UploadingList {
		if uploading.LocalFileMeta == nil {
			continue
		}
		if uploading.LocalFileMeta.EqualLengthMD5(meta) {
			return uploading.State
		}
		if uploading.LocalFileMeta.Path == meta.Path {
			// 移除旧的信息
			// 目前只是比较了文件大小
			if meta.Length != uploading.LocalFileMeta.Length {
				ud.Delete(meta)
				return nil
			}

			// 覆盖数据
			meta.MD5 = uploading.LocalFileMeta.MD5
			meta.SliceMD5 = uploading.LocalFileMeta.SliceMD5
			return uploading.State
		}
	}
	return nil
}

func (ud *UploadingDatabase) clearModTimeChange() {
	for i := 0; i < len(ud.UploadingList); i++ {
		uploading := ud.UploadingList[i]
		if uploading.LocalFileMeta == nil {
			continue
		}

		if uploading.ModTime == -1 { // 忽略
			continue
		}

		info, err := os.Stat(uploading.LocalFileMeta.Path)
		if err != nil {
			ud.deleteIndex(i)
			i--
			pcsUploadVerbose.Warnf("clear invalid file path: %s, err: %s\n", uploading.LocalFileMeta.Path, err)
			continue
		}

		if uploading.LocalFileMeta.ModTime != info.ModTime().Unix() {
			ud.deleteIndex(i)
			i--
			pcsUploadVerbose.Infof("clear modified file path: %s\n", uploading.LocalFileMeta.Path)
			continue
		}
	}
}

// Close 关闭数据库
func (ud *UploadingDatabase) Close() error {
	return ud.dataFile.Close()
}
