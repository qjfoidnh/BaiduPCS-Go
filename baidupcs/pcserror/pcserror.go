// Package pcserror PCS错误包
package pcserror

import (
	"github.com/iikira/BaiduPCS-Go/pcsutil/jsonhelper"
	"io"
)

type (
	// ErrType 错误类型
	ErrType int

	// Error 错误信息接口
	Error interface {
		error
		SetJSONError(err error)
		SetNetError(err error)
		SetRemoteError()
		GetOperation() string
		GetErrType() ErrType
		GetRemoteErrCode() int
		GetRemoteErrMsg() string
		GetError() error
	}
)

const (
	// ErrorTypeNoError 无错误
	ErrorTypeNoError ErrType = iota
	// ErrTypeInternalError 内部错误
	ErrTypeInternalError
	// ErrTypeRemoteError 远端服务器返回错误
	ErrTypeRemoteError
	// ErrTypeNetError 网络错误
	ErrTypeNetError
	// ErrTypeJSONParseError json 数据解析失败
	ErrTypeJSONParseError
	// ErrTypeOthers 其他错误
	ErrTypeOthers
)

const (
	// StrSuccess 操作成功
	StrSuccess = "操作成功"
	// StrInternalError 内部错误
	StrInternalError = "内部错误"
	// StrRemoteError 远端服务器返回错误
	StrRemoteError = "远端服务器返回错误"
	// StrNetError 网络错误
	StrNetError = "网络错误"
	// StrJSONParseError json 数据解析失败
	StrJSONParseError = "json 数据解析失败"
)

// DecodePCSJSONError 解析PCS JSON的错误
func DecodePCSJSONError(opreation string, data io.Reader) Error {
	errInfo := NewPCSErrorInfo(opreation)
	return HandleJSONParse(opreation, data, errInfo)
}

// DecodePanJSONError 解析Pan JSON的错误
func DecodePanJSONError(opreation string, data io.Reader) Error {
	errInfo := NewPanErrorInfo(opreation)
	return HandleJSONParse(opreation, data, errInfo)
}

// HandleJSONParse 处理解析json
func HandleJSONParse(op string, data io.Reader, info interface{}) (pcsError Error) {
	var (
		err     = jsonhelper.UnmarshalData(data, info)
		errInfo = info.(Error)
	)

	if errInfo == nil {
		errInfo = NewPCSErrorInfo(op)
	}

	if err != nil {
		errInfo.SetJSONError(err)
		return errInfo
	}

	// 设置出错类型为远程错误
	if errInfo.GetRemoteErrCode() != 0 {
		errInfo.SetRemoteError()
		return errInfo
	}

	return nil
}
