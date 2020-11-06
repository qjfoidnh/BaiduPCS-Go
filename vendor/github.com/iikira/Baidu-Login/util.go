package baidulogin

import (
	"fmt"
	"github.com/astaxie/beego/session"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

// registerBaiduClient 为 sess 如果没有 BaiduClient , 就添加
func registerBaiduClient(sess *session.Store) {
	if (*sess).Get("baiduclinet") == nil { // 找不到 cookie 储存器
		(*sess).Set("baiduclinet", NewBaiduClinet())
	}
}

// getBaiduClient 查找该 sessionID 下是否存在 BaiduClient
func getBaiduClient(sessionID string) (*BaiduClient, error) {
	sessionStore, err := globalSessions.GetSessionStore(sessionID)
	if err != nil {
		return NewBaiduClinet(), err
	}
	clientInterface := sessionStore.Get("baiduclinet")
	switch value := clientInterface.(type) {
	case *BaiduClient:
		return value, nil
	default:
		return NewBaiduClinet(), fmt.Errorf("Unknown session type: %s", value)
	}
}

// parseTemplate 自己写的简易 template 解析器
func parseTemplate(content string, rep map[string]string) string {
	for k, v := range rep {
		content = strings.Replace(content, "{{."+k+"}}", v, 1)
	}
	return content
}

// parsePhoneAndEmail 抓取绑定百度账号的邮箱和手机号并插入至 json 结构
func (lj *LoginJSON) parsePhoneAndEmail(bc *BaiduClient) {
	if lj.Data.GotoURL == "" {
		return
	}

	body, err := bc.Fetch("GET", lj.Data.GotoURL, nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	// 使用正则表达式匹配
	rawPhone := regexp.MustCompile(`<p class="verify-type-li-tiptop">(.*?)</p>\s+<p class="verify-type-li-tipbottom">通过手机验证码验证身份</p>`).FindSubmatch(body)
	rawEmail := regexp.MustCompile(`<p class="verify-type-li-tiptop">(.*?)</p>\s+<p class="verify-type-li-tipbottom">通过邮箱验证码验证身份</p>`).FindSubmatch(body)
	rawTokenAndU := regexp.MustCompile("token=(.*?)&u=(.*?)&secstate=").FindStringSubmatch(lj.Data.GotoURL)
	if len(rawPhone) >= 1 {
		lj.Data.Phone = string(rawPhone[1])
	} else {
		lj.Data.Phone = "未找到手机号"
	}

	if len(rawEmail) >= 1 {
		lj.Data.Email = string(rawEmail[1])
	} else {
		lj.Data.Email = "未找到邮箱地址"
	}

	if len(rawTokenAndU) >= 2 {
		lj.Data.Token = rawTokenAndU[1]
		if u, err := url.Parse(rawTokenAndU[2]); err == nil {
			lj.Data.U = u.Path
		}
	}
}

// parseCookies 解析 STOKEN, PTOKEN, BDUSS 并插入至 json 结构.
func (lj *LoginJSON) parseCookies(targetURL string, jar *cookiejar.Jar) {
	url, _ := url.Parse(targetURL)
	cookies := jar.Cookies(url)
	for _, cookie := range cookies {
		switch cookie.Name {
		case "BDUSS":
			lj.Data.BDUSS = cookie.Value
		case "PTOKEN":
			lj.Data.PToken = cookie.Value
		case "STOKEN":
			lj.Data.SToken = cookie.Value
		}
	}
	lj.Data.CookieString = pcsutil.GetURLCookieString(targetURL, jar) // 插入 cookie 字串
}
