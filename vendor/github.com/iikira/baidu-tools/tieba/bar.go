package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/baidu-tools/tieba/tiebautil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

// GetBars 获取贴吧列表
func (t *Tieba) GetBars() error {
	bars, err := GetBars(t.Baidu.UID)
	if err != nil {
		return err
	}
	t.Bars = bars
	return nil
}

// GetBars 通过 百度uid 获取贴吧列表
func GetBars(uid uint64) ([]*Bar, error) {
	var (
		pageNo uint16
		bars   []*Bar
	)
	bajsonRE := regexp.MustCompile("{\"id\":\".+?\"}")
	for {
		pageNo++
		rawQuery := fmt.Sprintf("_client_version=6.9.2.1&page_no=%d&page_size=200&uid=%d", pageNo, uid)

		//贴吧客户端签名
		body, err := requester.HTTPGet("http://c.tieba.baidu.com/c/f/forum/like?" + tiebautil.TiebaClientRawQuerySignature(rawQuery))
		if err != nil {
			return nil, fmt.Errorf("获取贴吧列表网络错误, %s", err)
		}

		if !strings.Contains(*(*string)(unsafe.Pointer(&body)), "has_more") { // 贴吧服务器响应有误, 再试一次
			pageNo--
			continue
		}

		jsonSlice := bajsonRE.FindAll(body, -1)
		if jsonSlice == nil { // 完成抓去贴吧列表
			break
		}

		for _, bajsonStr := range jsonSlice {
			bajson, err := simplejson.NewJson(bajsonStr)
			if err != nil {
				return nil, fmt.Errorf("获取贴吧列表json解析错误, %s", err)
			}
			if curScore, ok := bajson.CheckGet("cur_score"); ok {
				exp, _ := strconv.Atoi(curScore.MustString())
				bars = append(bars, &Bar{
					FID:   bajson.Get("id").MustString(),
					Name:  bajson.Get("name").MustString(),
					Level: bajson.Get("level_id").MustString(),
					Exp:   exp,
				})
			}
		}
	}
	return bars, nil
}

// GetTiebaFid 获取贴吧fid值
func GetTiebaFid(tiebaName string) (fid string, err error) {
	resp, err := requester.DefaultClient.Req("GET", "http://tieba.baidu.com/f/commit/share/fnameShareApi?ie=utf-8&fname="+tiebaName, nil, nil)
	if err != nil {
		return "", fmt.Errorf("获取贴吧fid网络错误, %s", err)
	}

	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("获取贴吧fid json解析错误, %s", err)
	}

	intFid := json.GetPath("data", "fid").MustInt()
	return fmt.Sprint(intFid), nil
}

// IsTiebaExist 检测贴吧是否存在
func IsTiebaExist(tiebaName string) bool {
	b, err := requester.HTTPGet("http://c.tieba.baidu.com/mo/q/m?tn4=bdKSW&sub4=&word=" + tiebaName)
	if err != nil {
		log.Println(err)
	}

	return !strings.Contains(*(*string)(unsafe.Pointer(&b)), `class="tip_text2">欢迎创建此吧，和朋友们在这里交流</p>`)
}
