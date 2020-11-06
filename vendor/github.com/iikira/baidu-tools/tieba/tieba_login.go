package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/baidu-tools"
	"github.com/iikira/baidu-tools/tieba/tiebautil"
	"strconv"
	"time"
)

// NewUserInfoByBDUSS 检测BDUSS有效性, 同时获取百度详细信息
func NewUserInfoByBDUSS(bduss string) (*Tieba, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	post := map[string]string{
		"bdusstoken":  bduss + "|null",
		"channel_id":  "",
		"channel_uid": "",
		"stErrorNums": "0",
		"subapp_type": "mini",
		"timestamp":   timestamp + "922",
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

	resp, err := requester.DefaultClient.Req("POST", "http://tieba.baidu.com/c/s/login", post, header) // 获取百度ID的TBS，UID，BDUSS等
	if err != nil {
		return nil, fmt.Errorf("检测BDUSS有效性网络错误, %s", err)
	}

	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("检测BDUSS有效性json解析出错: %s", err)
	}

	errCode := json.Get("error_code").MustString()
	errMsg := json.Get("error_msg").MustString()
	if errCode != "0" {
		return nil, fmt.Errorf("检测BDUSS有效性错误代码: %s, 消息: %s", errCode, errMsg)
	}

	userJSON := json.Get("user")
	uidStr := userJSON.Get("id").MustString()
	uid, _ := strconv.ParseUint(uidStr, 10, 64)

	t := &Tieba{
		Baidu: &baidu.Baidu{
			UID:  uid,
			Name: userJSON.Get("name").MustString(),
			Auth: &baidu.Auth{
				BDUSS: bduss,
			},
		},
		Tbs: json.GetPath("anti", "tbs").MustString(),
	}

	err = t.FlushUserInfo()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetTbs 获取贴吧TBS
func (t *Tieba) GetTbs() error {
	bduss := t.Baidu.Auth.BDUSS
	if bduss == "" {
		return fmt.Errorf("获取贴吧TBS出错: BDUSS为空")
	}
	tbs, err := GetTbs(bduss)
	if err != nil {
		return err
	}
	t.Tbs = tbs
	return nil
}

// GetTbs 通过 百度BDUSS 获取贴吧TBS
func GetTbs(bduss string) (tbs string, err error) {
	resp, err := requester.DefaultClient.Req("GET", "http://tieba.baidu.com/dc/common/tbs", nil, map[string]string{
		"Cookie": "BDUSS=" + bduss,
	})
	if err != nil {
		return "", fmt.Errorf("获取贴吧TBS网络错误: %s", err)
	}

	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("获取贴吧TBS json解析出错: %s", err)
	}

	isLogin := json.Get("is_login").MustInt()
	if isLogin != 0 {
		return json.Get("tbs").MustString(), nil
	}

	return "", fmt.Errorf("获取贴吧TBS错误, BDUSS无效")
}
