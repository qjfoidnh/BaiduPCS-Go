package baidupcs

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
)

type (
	// RecycleFDInfo 回收站中文件/目录信息
	RecycleFDInfo struct {
		FsID     int64  `json:"fs_id"` // fs_id
		Isdir    int    `json:"isdir"`
		LeftTime int    `json:"leftTime"`        // 剩余时间
		Path     string `json:"path"`            // 路径
		Filename string `json:"server_filename"` // 文件名 或 目录名
		Ctime    int64  `json:"server_ctime"`    // 创建日期
		Mtime    int64  `json:"server_mtime"`    // 修改日期
		MD5      string `json:"md5"`             // md5 值
		Size     int64  `json:"size"`            // 文件大小 (目录为0)
	}

	// RecycleFDInfoList 回收站中文件/目录列表
	RecycleFDInfoList []*RecycleFDInfo

	recycleListJSON struct {
		*pcserror.PanErrorInfo
		List RecycleFDInfoList `json:"list"`
	}

	recycleRestoreJSON struct {
		*pcserror.PCSErrInfo
		Extra FsIDListJSON `json:"extra"`
	}

	// RecycleClearInfo 清空回收站的信息
	RecycleClearInfo struct {
		List    RecycleFDInfoList `json:"list"`
		SussNum int               `json:"succNum"`
	}

	recycleClearJSON struct {
		*pcserror.PCSErrInfo
		Extra RecycleClearInfo `json:"extra"`
	}
)

// RecycleList 列出回收站文件列表
func (pcs *BaiduPCS) RecycleList(page int) (fdl RecycleFDInfoList, panError pcserror.Error) {
	dataReadCloser, panError := pcs.PrepareRecycleList(page)
	if panError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationRecycleList)
	jsonData := recycleListJSON{
		PanErrorInfo: errInfo,
	}

	panError = pcserror.HandleJSONParse(OperationRecycleList, dataReadCloser, &jsonData)
	if panError != nil {
		return
	}

	return jsonData.List, nil
}

// RecycleRestore 还原回收站文件或目录
func (pcs *BaiduPCS) RecycleRestore(fidList ...int64) (sussFsIDList []*FsIDJSON, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareRecycleRestore(fidList...)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPCSErrorInfo(OperationRecycleRestore)
	jsonData := recycleRestoreJSON{
		PCSErrInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationRecycleRestore, dataReadCloser, &jsonData)
	return jsonData.Extra.List, pcsError
}

// RecycleDelete 删除回收站文件或目录
func (pcs *BaiduPCS) RecycleDelete(fidList ...int64) (panError pcserror.Error) {
	dataReadCloser, panError := pcs.PrepareRecycleDelete(fidList...)
	if panError != nil {
		return
	}

	defer dataReadCloser.Close()

	panError = pcserror.DecodePanJSONError(OperationRecycleDelete, dataReadCloser)
	return
}

// RecycleClear 清空回收站
func (pcs *BaiduPCS) RecycleClear() (sussNum int, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareRecycleClear()
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPCSErrorInfo(OperationRecycleClear)
	jsonData := recycleClearJSON{
		PCSErrInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationRecycleClear, dataReadCloser, &jsonData)
	return jsonData.Extra.SussNum, pcsError
}
