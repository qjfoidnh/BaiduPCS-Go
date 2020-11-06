package pcserror

import (
	"fmt"
)

type (
	// PanErrorInfo 网盘网页的api错误
	PanErrorInfo struct {
		Operation string
		ErrType   ErrType
		Err       error
		ErrNo     int `json:"errno"`
		// ErrMsg    string `json:"err_msg"`
	}
)

// NewPanErrorInfo 提供operation操作名称, 返回 *PanErrorInfo
func NewPanErrorInfo(operation string) *PanErrorInfo {
	return &PanErrorInfo{
		Operation: operation,
		ErrType:   ErrorTypeNoError,
	}
}

// SetJSONError 设置JSON错误
func (pane *PanErrorInfo) SetJSONError(err error) {
	pane.ErrType = ErrTypeJSONParseError
	pane.Err = err
}

// SetNetError 设置网络错误
func (pane *PanErrorInfo) SetNetError(err error) {
	pane.ErrType = ErrTypeNetError
	pane.Err = err
}

// SetRemoteError 设置远端服务器错误
func (pane *PanErrorInfo) SetRemoteError() {
	pane.ErrType = ErrTypeRemoteError
}

// GetOperation 获取操作
func (pane *PanErrorInfo) GetOperation() string {
	return pane.Operation
}

// GetErrType 获取错误类型
func (pane *PanErrorInfo) GetErrType() ErrType {
	return pane.ErrType
}

// GetRemoteErrCode 获取远端服务器错误代码
func (pane *PanErrorInfo) GetRemoteErrCode() int {
	return pane.ErrNo
}

// GetRemoteErrMsg 获取远端服务器错误消息
func (pane *PanErrorInfo) GetRemoteErrMsg() string {
	return FindPanErr(pane.ErrNo)
}

// GetError 获取原始错误
func (pane *PanErrorInfo) GetError() error {
	return pane.Err
}

func (pane *PanErrorInfo) Error() string {
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

		errmsg := FindPanErr(pane.ErrNo)
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
func FindPanErr(errno int) (errmsg string) {
	switch errno {
	case 0:
		return StrSuccess
	case -1:
		return "由于您分享了违反相关法律法规的文件，分享功能已被禁用，之前分享出去的文件不受影响。"
	case -2:
		return "用户不存在,请刷新页面后重试"
	case -3:
		return "文件不存在,请刷新页面后重试"
	case -4:
		return "登录信息有误，请重新登录试试"
	case -5:
		return "host_key和user_key无效"
	case -6:
		return "请重新登录"
	case -7:
		return "该分享已删除或已取消"
	case -8:
		return "该分享已经过期"
	case -9:
		return "文件不存在"
	case -10:
		return "分享外链已经达到最大上限100000条，不能再次分享"
	case -11:
		return "验证cookie无效"
	case -12:
		return "访问密码错误"
	case -14:
		return "对不起，短信分享每天限制20条，你今天已经分享完，请明天再来分享吧！"
	case -15:
		return "对不起，邮件分享每天限制20封，你今天已经分享完，请明天再来分享吧！"
	case -16:
		return "对不起，该文件已经限制分享！"
	case -17:
		return "文件分享超过限制"
	case -19:
		return "需要输入验证码"
	case -21:
		return "分享已取消或分享信息无效"
	case -30:
		return "文件已存在"
	case -31:
		return "文件保存失败"
	case -33:
		return "一次支持操作999个，减点试试吧"
	case -62:
		return "可能需要输入验证码"
	case -70:
		return "你分享的文件中包含病毒或疑似病毒，为了你和他人的数据安全，换个文件分享吧"
	case 2:
		return "参数错误"
	case 3:
		return "未登录或帐号无效"
	case 4:
		return "存储好像出问题了，请稍候再试"
	case 105:
		return "啊哦，链接错误没找到文件，请打开正确的分享链接"
	case 108:
		return "文件名有敏感词，优化一下吧"
	case 110:
		return "分享次数超出限制，可以到“我的分享”中查看已分享的文件链接"
	case 112:
		return "页面已过期，请刷新后重试"
	case 113:
		return "签名错误"
	case 114:
		return "当前任务不存在，保存失败"
	case 115:
		return "该文件禁止分享"
	case 132:
		return "您的帐号可能存在安全风险，为了确保为您本人操作，请先进行安全验证。"
	default:
		return "未知错误"
	}
}
