package pcserror

import (
	"fmt"
)

type (
	// PanErrorInfo 网盘网页的api错误
	XPanErrorInfo struct {
		Operation string
		ErrType   ErrType
		Err       error
		ErrNo     int `json:"errno"`
	}
)

// NewXPanErrorInfo 提供operation操作名称, 返回 *PanErrorInfo
func NewXPanErrorInfo(operation string) *XPanErrorInfo {
	return &XPanErrorInfo{
		Operation: operation,
		ErrType:   ErrorTypeNoError,
	}
}

// SetJSONError 设置JSON错误
func (pane *XPanErrorInfo) SetJSONError(err error) {
	pane.ErrType = ErrTypeJSONParseError
	pane.Err = err
}

// SetNetError 设置网络错误
func (pane *XPanErrorInfo) SetNetError(err error) {
	pane.ErrType = ErrTypeNetError
	pane.Err = err
}

// SetRemoteError 设置远端服务器错误
func (pane *XPanErrorInfo) SetRemoteError() {
	pane.ErrType = ErrTypeRemoteError
}

// GetOperation 获取操作
func (pane *XPanErrorInfo) GetOperation() string {
	return pane.Operation
}

// GetErrType 获取错误类型
func (pane *XPanErrorInfo) GetErrType() ErrType {
	return pane.ErrType
}

// GetRemoteErrCode 获取远端服务器错误代码
func (pane *XPanErrorInfo) GetRemoteErrCode() int {
	return pane.ErrNo
}

// GetRemoteErrMsg 获取远端服务器错误消息
func (pane *XPanErrorInfo) GetRemoteErrMsg() string {
	return FindPanErr(pane.ErrNo)
}

// GetError 获取原始错误
func (pane *XPanErrorInfo) GetError() error {
	return pane.Err
}

func (pane *XPanErrorInfo) Error() string {
	if pane.Operation == "" {
		if pane.Err != nil {
			return pane.Err.Error()
		}
		return StrSuccess
	}

	switch pane.ErrType {
	case ErrTypeInternalError:
		return fmt.Sprintf("%s: %s, %s", pane.Operation, StrInternalError, pane.Err)
	case ErrTypeJSONParseError:
		return fmt.Sprintf("%s: %s, %s", pane.Operation, StrJSONParseError, pane.Err)
	case ErrTypeNetError:
		return fmt.Sprintf("%s: %s, %s", pane.Operation, StrNetError, pane.Err)
	case ErrTypeRemoteError:
		if pane.ErrNo == 0 {
			return fmt.Sprintf("%s: %s", pane.Operation, StrSuccess)
		}

		errmsg := FindXPanErr(pane.ErrNo)
		return fmt.Sprintf("%s: 遇到错误, %s, 代码: %d, 消息: %s", pane.Operation, StrRemoteError, pane.ErrNo, errmsg)
	case ErrTypeOthers:
		if pane.Err == nil {
			return fmt.Sprintf("%s: %s", pane.Operation, StrSuccess)
		}

		return fmt.Sprintf("%s, 遇到错误, %s", pane.Operation, pane.Err)
	default:
		panic("panerrorinfo: unknown ErrType")
	}
}

// FindPanErr 根据 ErrNo, 解析网盘错误信息
func FindXPanErr(errno int) (errmsg string) {
	switch errno {
	case 0:
		return StrSuccess
	case 2:
		return "md5不存在"
	case -8:
		return "同名文件冲突"
	default:
		return "未知错误"
	}
}
