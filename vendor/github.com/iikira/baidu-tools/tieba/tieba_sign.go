package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/baidu-tools/tieba/tiebautil"
	"strconv"
	"time"
	"unsafe"
)

// TiebaSign 贴吧签到
func (user *Tieba) TiebaSign(fid, name string) (errorCode, errorMsg string, bonusExp int, err error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	post := map[string]string{
		"BDUSS":       user.Baidu.Auth.BDUSS,
		"_client_id":  "wappc_" + timestamp + "150_607",
		"fid":         fid,
		"kw":          name,
		"stErrorNums": "1",
		"stMethod":    "1",
		"stMode":      "1",
		"stSize":      "229",
		"stTime":      "185",
		"stTimesNum":  "1",
		"subapp_type": "mini",
		"tbs":         user.Tbs,
		"timestamp":   timestamp + "083",
	}
	tiebautil.TiebaClientSignature(post)

	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Cookie":       "ka=open",
		"net":          "1",
		"User-Agent":   "bdtb for Android 6.9.2.1",
		"client_logid": timestamp + "416",
		"Connection":   "Keep-Alive",
	}

	body, err := requester.DefaultClient.Fetch("POST", "http://c.tieba.baidu.com/c/c/forum/sign", post, header)
	if err != nil {
		return "", "", 0, fmt.Errorf("贴吧签到网络错误: %s", err)
	}

	json, err := simplejson.NewJson(body)
	if err != nil {
		return "", "", 0, fmt.Errorf("贴吧签到json解析错误: %s", err)
	}

	errorCode = json.Get("error_code").MustString()
	errorMsg = json.Get("error_msg").MustString()

	if signBonusPoint, ok := json.Get("user_info").CheckGet("sign_bonus_point"); ok { // 签到成功, 获取经验
		bonusExp, _ = strconv.Atoi(signBonusPoint.MustString())
		return errorCode, errorMsg, bonusExp, nil
	}

	if errorMsg == "" {
		return errorCode, errorMsg, 0, fmt.Errorf("贴吧签到时发生错误, 未能找到错误原因, 请检查：" + *(*string)(unsafe.Pointer(&body)))
	}

	return errorCode, errorMsg, 0, nil
}

// DoTiebaSign 执行贴吧签到
func (user *Tieba) DoTiebaSign(fid, name string) (status int, bonusExp int, err error) {
	errorCode, errorMsg, bonusExp, err := user.TiebaSign(fid, name)
	if err != nil {
		return 1, bonusExp, err
	}

	err = fmt.Errorf("贴吧签到时发生错误, 错误代码: %s, 消息: %s", pcsutil.ErrorColor(errorCode), pcsutil.ErrorColor(errorMsg))
	switch errorCode {
	case "0", "160002": // 	签到成功 / 已签到
		return 0, bonusExp, nil
	case "220034", "340011": // 操作太快
		return 2, bonusExp, err
	case "300000": // 未知错误
		return 2, bonusExp, err
	case "340008", "340006", "110001", "3250002": // 340008黑名单, 340006封吧, 110001签名错误, 3250002永久封号
		return 3, bonusExp, err
	case "1", "1990055": // 1掉线, 1990055未实名
		return 4, bonusExp, err
	default:
		return 1, bonusExp, err
	}
}
