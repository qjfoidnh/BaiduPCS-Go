package baidulogin

import (
	"bytes"
	"fmt"
	"github.com/iikira/Baidu-Login/bdcrypto"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/json-iterator/go"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"time"
)

// BaiduClient 记录登录百度所使用的信息
type BaiduClient struct {
	*requester.HTTPClient

	serverTime          string // 百度服务器时间, 形如 "e362bacbae"
	rsaPublicKeyModulus string
	fpUID               string
	traceid             string
}

// LoginJSON 从百度服务器解析的数据结构
type LoginJSON struct {
	ErrInfo struct {
		No  string `json:"no"`
		Msg string `json:"msg"`
	} `json:"errInfo"`
	Data struct {
		CodeString   string `json:"codeString"`
		GotoURL      string `json:"gotoUrl"`
		Token        string `json:"token"`
		U            string `json:"u"`
		AuthSID      string `json:"authsid"`
		Phone        string `json:"phone"`
		Email        string `json:"email"`
		BDUSS        string `json:"bduss"`
		PToken       string `json:"ptoken"`
		SToken       string `json:"stoken"`
		CookieString string `json:"cookieString"`
	} `json:"data"`
}

// NewBaiduClinet 返回 BaiduClient 指针对象
func NewBaiduClinet() *BaiduClient {
	bc := &BaiduClient{
		HTTPClient: requester.NewHTTPClient(),
	}

	bc.getServerTime()               // 访问一次百度页面，以初始化百度的 Cookie
	bc.getBaiduRSAPublicKeyModulus() //
	bc.getTraceID()
	return bc
}

// BaiduLogin 发送 百度登录请求
func (bc *BaiduClient) BaiduLogin(username, password, verifycode, vcodestr string) (lj *LoginJSON) {
	lj = &LoginJSON{}
	enpass, err := bdcrypto.RSAEncryptOfWapBaidu(bc.rsaPublicKeyModulus, []byte(password+bc.serverTime))
	if err != nil {
		lj.ErrInfo.No = "-1"
		lj.ErrInfo.Msg = "RSA加密失败, " + err.Error() + ": " + bc.rsaPublicKeyModulus
		return lj
	}

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10) + "773_357"
	post := map[string]string{
		"username":     username,
		"password":     enpass,
		"verifycode":   verifycode,
		"vcodestr":     vcodestr,
		"isphone":      "0",
		"loginmerge":   "1", // 加入这个, 不用判断是否手机号了
		"action":       "login",
		"uid":          timestampStr,
		"skin":         "default_v2",
		"connect":      "0",
		"dv":           "tk0.408376350146535171516806245342@oov0QqrkqfOuwaCIxUELn3oYlSOI8f51tbnGy-nk3crkqfOuwaCIxUou2iobENoYBf51tb4Gy-nk3cuv0ounk5vrkBynGyvn1QzruvN6z3drLJi6LsdFIe3rkt~4Lyz5ktfn1Qlrk5v5D5fOuwaCIxUobJWOI3~rkt~4Lyi5kBfni0vrk8~n15fOuwaCIxUobJWOI3~rkt~4Lyz5DQfn1oxrk0v5k5eruvN6z3drLneFYeVEmy-nk3c-qq6Cqw3h7CChwvi5-y-rkFizvmEufyr1By4k5bn15e5k0~n18inD0b5D8vn1Tyn1t~nD5~5T__ivmCpA~op5gr-wbFLhyFLnirYsSCIAerYnNOGcfEIlQ6I6VOYJQIvh515f51tf5DBv5-yln15f5DFy5myl5kqf5DFy5myvnktxrkT-5T__Hv0nq5myv5myv4my-nWy-4my~n-yz5myz4Gyx4myv5k0f5Dqirk0ynWyv5iTf5DB~rk0z5Gyv4kTf5DQxrkty5Gy-5iQf51B-rkt~4B__",
		"getpassUrl":   "/passport/getpass?clientfrom=&adapter=0&ssid=&from=&authsite=&bd_page_type=&uid=" + timestampStr + "&pu=&tpl=wimn&u=https://m.baidu.com/usrprofile%3Fuid%3D" + timestampStr + "%23logined&type=&bdcm=060d5ffd462309f7e5529822720e0cf3d7cad665&tn=&regist_mode=&login_share_strategy=&subpro=wimn&skin=default_v2&client=&connect=0&smsLoginLink=1&loginLink=&bindToSmsLogin=&overseas=&is_voice_sms=&subpro=wimn&hideSLogin=&forcesetpwd=&regdomestic=",
		"mobilenum":    "undefined",
		"servertime":   bc.serverTime,
		"gid":          "DA7C3AE-AF1F-48C0-AF9C-F1882CA37CD5",
		"logLoginType": "wap_loginTouch",
		"FP_UID":       "0b58c206c9faa8349576163341ef1321",
		"traceid":      bc.traceid,
	}

	header := map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded",
		"Accept":           "application/json",
		"Referer":          "https://wappass.baidu.com/",
		"X-Requested-With": "XMLHttpRequest",
		"Connection":       "keep-alive",
	}

	body, err := bc.Fetch("POST", "https://wappass.baidu.com/wp/api/login", post, header)
	if err != nil {
		lj.ErrInfo.No = "-1"
		lj.ErrInfo.Msg = "网络请求失败, " + err.Error()
		return lj
	}

	// 如果 json 解析出错
	if err = jsoniter.Unmarshal(body, &lj); err != nil {
		lj.ErrInfo.No = "-1"
		lj.ErrInfo.Msg = "发送登录请求错误: " + err.Error()
		return lj
	}

	switch lj.ErrInfo.No {
	case "0":
		lj.parseCookies("https://wappass.baidu.com", bc.Jar.(*cookiejar.Jar)) // 解析登录数据
	case "400023", "400101": // 需要验证手机或邮箱
		lj.parsePhoneAndEmail(bc)
	case "400408": // 应国家相关法律要求，自6月1日起使用信息发布、即时通讯等互联网服务需进行身份信息验证。为保障您对相关服务功能的正常使用，建议您尽快完成手机号验证，感谢您的理解和支持。
	}

	return lj
}

// SendCodeToUser 发送验证码到 手机/邮箱
func (bc *BaiduClient) SendCodeToUser(verifyType, token string) (msg string) {
	url := fmt.Sprintf("https://wappass.baidu.com/passport/authwidget?action=send&tpl=&type=%s&token=%s&from=&skin=&clientfrom=&adapter=2&updatessn=&bindToSmsLogin=&upsms=&finance=", verifyType, token)
	body, err := bc.Fetch("GET", url, nil, nil)
	if err != nil {
		return err.Error()
	}

	rawMsg := regexp.MustCompile(`<p class="mod-tipinfo-subtitle">\s+(.*?)\s+</p>`).FindSubmatch(body)
	if len(rawMsg) >= 1 {
		return string(rawMsg[1])
	}

	return "未知消息"
}

// VerifyCode 输入 手机/邮箱 收到的验证码, 验证登录
func (bc *BaiduClient) VerifyCode(verifyType, token, vcode, u string) (lj *LoginJSON) {
	lj = &LoginJSON{}
	header := map[string]string{
		"Connection":                "keep-alive",
		"Host":                      "wappass.baidu.com",
		"Pragma":                    "no-cache",
		"Upgrade-Insecure-Requests": "1",
	}

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10) + "773_357" + "994"
	url := fmt.Sprintf("https://wappass.baidu.com/passport/authwidget?v="+timestampStr+"&vcode=%s&token=%s&u=%s&action=check&type=%s&tpl=&skin=&clientfrom=&adapter=2&updatessn=&bindToSmsLogin=&isnew=&card_no=&finance=&callback=%s", vcode, token, u, verifyType, "jsonp1")
	body, err := bc.Fetch("GET", url, nil, header)
	if err != nil {
		lj.ErrInfo.No = "-2"
		lj.ErrInfo.Msg = "网络请求错误: " + err.Error()
		return
	}

	// 去除 body 的 callback 嵌套 "jsonp1(...)"
	body = bytes.TrimPrefix(body, []byte("jsonp1("))
	body = bytes.TrimSuffix(body, []byte(")"))

	// 如果 json 解析出错, 直接输出错误信息
	if err := jsoniter.Unmarshal(body, &lj); err != nil {
		lj.ErrInfo.No = "-2"
		lj.ErrInfo.Msg = "提交手机/邮箱验证码错误: " + err.Error()
		return
	}

	// 最后一步要访问的 URL
	u = fmt.Sprintf("%s&authsid=%s&fromtype=%s&bindToSmsLogin=", u, lj.Data.AuthSID, verifyType) // url

	_, err = bc.Fetch("GET", u, nil, nil)
	if err != nil {
		lj.ErrInfo.No = "-2"
		lj.ErrInfo.Msg = "提交手机/邮箱验证码错误: " + err.Error()
		return
	}

	lj.parseCookies(u, bc.Jar.(*cookiejar.Jar))
	return lj
}

// getTraceID 获取百度 Trace-Id
func (bc *BaiduClient) getTraceID() {
	resp, err := bc.Req("GET", "https://wappass.baidu.com/", nil, nil)
	if err != nil {
		bc.traceid = err.Error()
		return
	}
	bc.traceid = resp.Header.Get("Trace-Id")
	resp.Body.Close()
}

// getServerTime 获取百度服务器时间, 形如 "e362bacbae"
func (bc *BaiduClient) getServerTime() {
	body, _ := bc.Fetch("GET", "https://wappass.baidu.com/wp/api/security/antireplaytoken", nil, nil)
	rawServerTime := regexp.MustCompile(`,"time":"(.*?)"`).FindSubmatch(body)
	if len(rawServerTime) >= 1 {
		bc.serverTime = string(rawServerTime[1])
		return
	}
	bc.serverTime = "e362bacbae"
}

// getBaiduRSAPublicKeyModulus 获取百度 RSA 字串
func (bc *BaiduClient) getBaiduRSAPublicKeyModulus() {
	body, _ := bc.Fetch("GET", "https://wappass.baidu.com/static/touch/js/login_d9bffc9.js", nil, nil)
	rawRSA := regexp.MustCompile(`,rsa:"(.*?)",error:`).FindSubmatch(body)
	if len(rawRSA) >= 1 {
		bc.rsaPublicKeyModulus = string(rawRSA[1])
		return
	}
	bc.rsaPublicKeyModulus = "B3C61EBBA4659C4CE3639287EE871F1F48F7930EA977991C7AFE3CC442FEA49643212E7D570C853F368065CC57A2014666DA8AE7D493FD47D171C0D894EEE3ED7F99F6798B7FFD7B5873227038AD23E3197631A8CB642213B9F27D4901AB0D92BFA27542AE890855396ED92775255C977F5C302F1E7ED4B1E369C12CB6B1822F"
}
