package pcserror

import (
	"fmt"
)

type (
	// PCSErrInfo PCS错误信息
	PCSErrInfo struct {
		Operation string // 正在进行的操作
		ErrType   ErrType
		Err       error
		ErrCode   int    `json:"error_code"` // 错误代码
		ErrMsg    string `json:"error_msg"`  // 错误消息
	}
)

// NewPCSErrorInfo 提供operation操作名称, 返回 *PCSErrInfo
func NewPCSErrorInfo(operation string) *PCSErrInfo {
	return &PCSErrInfo{
		Operation: operation,
		ErrType:   ErrorTypeNoError,
	}
}

// SetJSONError 设置JSON错误
func (pcse *PCSErrInfo) SetJSONError(err error) {
	pcse.ErrType = ErrTypeJSONParseError
	pcse.Err = err
}

// SetNetError 设置网络错误
func (pcse *PCSErrInfo) SetNetError(err error) {
	pcse.ErrType = ErrTypeNetError
	pcse.Err = err
}

// SetRemoteError 设置远端服务器错误
func (pcse *PCSErrInfo) SetRemoteError() {
	pcse.ErrType = ErrTypeRemoteError
}

// GetOperation 获取操作
func (pcse *PCSErrInfo) GetOperation() string {
	return pcse.Operation
}

// GetErrType 获取错误类型
func (pcse *PCSErrInfo) GetErrType() ErrType {
	return pcse.ErrType
}

// GetRemoteErrCode 获取远端服务器错误代码
func (pcse *PCSErrInfo) GetRemoteErrCode() int {
	return pcse.ErrCode
}

// GetRemoteErrMsg 获取远端服务器错误消息
func (pcse *PCSErrInfo) GetRemoteErrMsg() string {
	_, msg := findPCSErr(pcse.ErrCode, pcse.ErrMsg)
	return msg
}

// GetError 获取原始错误
func (pcse *PCSErrInfo) GetError() error {
	return pcse.Err
}

func (pcse *PCSErrInfo) Error() string {
	if pcse.Operation == "" {
		if pcse.Err != nil {
			return pcse.Err.Error()
		}
		return StrSuccess
	}

	switch pcse.ErrType {
	case ErrTypeInternalError:
		return fmt.Sprintf("%s: %s, %s", pcse.Operation, StrInternalError, pcse.Err)
	case ErrTypeJSONParseError:
		return fmt.Sprintf("%s: %s, %s", pcse.Operation, StrJSONParseError, pcse.Err)
	case ErrTypeNetError:
		return fmt.Sprintf("%s: %s, %s", pcse.Operation, StrNetError, pcse.Err)
	case ErrTypeRemoteError:
		if pcse.ErrCode == 0 {
			return fmt.Sprintf("%s: %s", pcse.Operation, StrSuccess)
		}

		code, msg := findPCSErr(pcse.ErrCode, pcse.ErrMsg)
		return fmt.Sprintf("%s: 遇到错误, %s, 代码: %d, 消息: %s", pcse.Operation, StrRemoteError, code, msg)
	case ErrTypeOthers:
		if pcse.Err == nil {
			return fmt.Sprintf("%s: %s", pcse.Operation, StrSuccess)
		}

		return fmt.Sprintf("%s, 遇到错误, %s", pcse.Operation, pcse.Err)
	default:
		panic("pcserrorinfo: unknown ErrType")
	}
}

// findPCSErr 检查 PCS 错误, 查找已知错误
func findPCSErr(errCode int, errMsg string) (int, string) {
	switch errCode {
	case 0:
		return errCode, ""
	case 31045: // user not exists
		return errCode, "操作失败, 可能百度帐号登录状态过期, 请尝试重新登录, 消息: " + errMsg
	case 31066: // file does not exist
		return errCode, "文件或目录不存在"
	case 31079: // file md5 not found, you should use upload api to upload the whole file.
		return errCode, "秒传文件失败"
	}
	return errCode, errMsg
}
