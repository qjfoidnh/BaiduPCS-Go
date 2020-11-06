package baidupcs

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
)

type quotaInfo struct {
	*pcserror.PCSErrInfo
	Quota int64 `json:"quota"`
	Used  int64 `json:"used"`
}

// QuotaInfo 获取当前用户空间配额信息
func (pcs *BaiduPCS) QuotaInfo() (quota, used int64, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareQuotaInfo()
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	quotaInfo := &quotaInfo{
		PCSErrInfo: pcserror.NewPCSErrorInfo(OperationQuotaInfo),
	}

	pcsError = pcserror.HandleJSONParse(OperationQuotaInfo, dataReadCloser, quotaInfo)
	if pcsError != nil {
		return
	}

	return quotaInfo.Quota, quotaInfo.Used, nil
}
