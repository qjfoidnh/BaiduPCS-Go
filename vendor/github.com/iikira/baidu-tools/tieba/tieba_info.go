package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/baidu-tools"
	"github.com/iikira/baidu-tools/tieba/tiebautil"
	"strconv"
)

// NewUserInfoByUID 提供 UID 获取百度帐号详细信息
func NewUserInfoByUID(uid uint64) (t *Tieba, err error) {
	b := &baidu.Baidu{
		UID: uid,
	}

	rawQuery := "has_plist=0&need_post_count=1&rn=1&uid=" + fmt.Sprint(b.UID)
	urlStr := "http://c.tieba.baidu.com/c/u/user/profile?" + tiebautil.TiebaClientRawQuerySignature(rawQuery)
	resp, err := requester.DefaultClient.Req("GET", urlStr, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	userJSON := json.GetPath("user")
	b.Name = userJSON.Get("name").MustString()
	b.NameShow = userJSON.Get("name_show").MustString()

	ageStr := userJSON.Get("tb_age").MustString()
	b.Age, _ = strconv.ParseFloat(ageStr, 64)

	sex := userJSON.Get("sex").MustInt()
	switch sex {
	case 1:
		b.Sex = "♂"
	case 2:
		b.Sex = "♀"
	default:
		b.Sex = "unknown"
	}

	t = &Tieba{
		Baidu: b,
		Stat: &Stat{
			LikeForumNum: userJSON.Get("like_forum_num").MustInt(),
			PostNum:      userJSON.Get("post_num").MustInt(),
		},
	}

	return t, nil
}

// NewUserInfoByName 提供 name (百度用户名) 获取百度帐号个人详细信息
func NewUserInfoByName(name string) (t *Tieba, err error) {
	resp, err := requester.DefaultClient.Req("GET", "http://tieba.baidu.com/home/get/panel?un="+name, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return NewUserInfoByUID(json.GetPath("data", "id").MustUint64())
}

// FlushUserInfo 刷新百度帐号个人详细信息
func (t *Tieba) FlushUserInfo() (err error) {
	if t.Baidu == nil {
		return fmt.Errorf("Baidu is not initialize")
	}

	this := new(Tieba)

	if t.Baidu.UID != 0 {
		this, err = NewUserInfoByUID(t.Baidu.UID)
		if err != nil {
			return err
		}
	} else if t.Baidu.Name != "" {
		this, err = NewUserInfoByName(t.Baidu.Name)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Baidu uid and name are null")
	}

	this.Baidu.Auth = t.Baidu.Auth
	t.Baidu = this.Baidu
	t.Stat = this.Stat
	return nil
}
