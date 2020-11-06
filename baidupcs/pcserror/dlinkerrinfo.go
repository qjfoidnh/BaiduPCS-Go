package pcserror

import (
	"fmt"
)

type (
	// DlinkErrInfo dlink服务器错误信息
	DlinkErrInfo struct {
		Operation string
		ErrType   ErrType
		Err       error
		ErrNo     int    `json:"errno"`
		Msg       string `json:"msg"`
	}
)

// NewDlinkErrInfo 初始化DlinkErrInfo
func NewDlinkErrInfo(op string) *DlinkErrInfo {
	return &DlinkErrInfo{
		Operation: op,
	}
}

// SetJSONError 设置JSON错误
func (dle *DlinkErrInfo) SetJSONError(err error) {
	dle.ErrType = ErrTypeJSONParseError
	dle.Err = err
}

// SetNetError 设置网络错误
func (dle *DlinkErrInfo) SetNetError(err error) {
	dle.ErrType = ErrTypeNetError
	dle.Err = err
}

// SetRemoteError 设置远端服务器错误
func (dle *DlinkErrInfo) SetRemoteError() {
	dle.ErrType = ErrTypeRemoteError
}

// GetOperation 获取操作
func (dle *DlinkErrInfo) GetOperation() string {
	return dle.Operation
}

// GetErrType 获取错误类型
func (dle *DlinkErrInfo) GetErrType() ErrType {
	return dle.ErrType
}

// GetRemoteErrCode 获取远端服务器错误代码
func (dle *DlinkErrInfo) GetRemoteErrCode() int {
	return dle.ErrNo
}

// GetRemoteErrMsg 获取远端服务器错误消息
func (dle *DlinkErrInfo) GetRemoteErrMsg() string {
	return dle.Msg
}

// GetError 获取原始错误
func (dle *DlinkErrInfo) GetError() error {
	return dle.Err
}

func (dle *DlinkErrInfo) Error() string {
	if dle.Operation == "" {
		if dle.Err != nil {
			return dle.Err.Error()
		}
		return StrSuccess
	}

	switch dle.ErrType {
	case ErrTypeInternalError:
		return fmt.Sprintf("%s: %s, %s", dle.Operation, StrInternalError, dle.Err)
	case ErrTypeJSONParseError:
		return fmt.Sprintf("%s: %s, %s", dle.Operation, StrJSONParseError, dle.Err)
	case ErrTypeNetError:
		return fmt.Sprintf("%s: %s, %s", dle.Operation, StrNetError, dle.Err)
	case ErrTypeRemoteError:
		if dle.ErrNo == 0 {
			return fmt.Sprintf("%s: %s", dle.Operation, StrSuccess)
		}

		return fmt.Sprintf("%s: 遇到错误, %s, 代码: %d, 消息: %s", dle.Operation, StrRemoteError, dle.ErrNo, dle.Msg)
	case ErrTypeOthers:
		if dle.Err == nil {
			return fmt.Sprintf("%s: %s", dle.Operation, StrSuccess)
		}

		return fmt.Sprintf("%s, 遇到错误, %s", dle.Operation, dle.Err)
	default:
		panic("dlinkerrinfo: unknown ErrType")
	}
}
