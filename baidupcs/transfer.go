package baidupcs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

func (pcs *BaiduPCS) GenerateShareQueryURL(subPath string, params map[string]string) *url.URL {
	shareURL := &url.URL{
		Scheme: GetHTTPScheme(true),
		Host:   PanBaiduCom,
		Path:   "/share/" + subPath,
	}
	uv := shareURL.Query()
	uv.Set("app_id", PanAppID)
	uv.Set("channel", "chunlei")
	uv.Set("clienttype", "0")
	uv.Set("web", "1")
	for key, value := range params {
		uv.Set(key, value)
	}

	shareURL.RawQuery = uv.Encode()
	return shareURL
}

func (pcs *BaiduPCS) PostShareQuery(url string, data map[string]string) (res map[string]string) {
	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodPost, url, data, map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      "https://pan.baidu.com/disk/home",
	})
	res = make(map[string]string)
	if panError != nil {
		res["ErrMsg"] = "提交分享项查询请求时发生错误"
		return
	}
	defer dataReadCloser.Close()
	body, _ := ioutil.ReadAll(dataReadCloser)
	errno := gjson.Get(string(body), `errno`).String()
	if errno != "0" {
		res["ErrMsg"] = "分享码错误"
		return
	}
	res["ErrMsg"] = "0"
	return
}

func (pcs *BaiduPCS) AccessSharePage(sharelink string, first bool) (tokens map[string]string) {
	tokens = make(map[string]string)
	tokens["ErrMsg"] = "0"
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
	headers["Referer"] = "https://pan.baidu.com/disk/home"

	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodGet, sharelink, nil, headers)

	if panError != nil {
		tokens["ErrMsg"] = "访问分享页失败"
		return
	}
	defer dataReadCloser.Close()
	body, _ := ioutil.ReadAll(dataReadCloser)
	not_found_flag := strings.Contains(string(body), "platform-non-found")
	error_page_title := strings.Contains(string(body), "error-404")
	if error_page_title {
		tokens["ErrMsg"] = "页面不存在"
		return
	}
	if not_found_flag {
		tokens["ErrMsg"] = "分享链接已失效"
		return
	} else {
		re, _ := regexp.Compile(`"bdstoken":"([A-Za-z0-9]+)",`)
		sub := re.FindSubmatch(body)
		if len(sub) < 2 {
			tokens["ErrMsg"] = "分享页面解析失败"
			return
		}
		tokens["bdstoken"] = string(sub[1])
		if first {
			return
		}
		re, _ = regexp.Compile(`"shareid":([0-9]+),`)
		sub = re.FindSubmatch(body)
		if len(sub) < 2 && !first {
			tokens["ErrMsg"] = "缺少提取码"
			return
		}
		tokens["shareid"] = string(sub[1])
		re, _ = regexp.Compile(`uk=([0-9]+)`)
		sub = re.FindSubmatch(body)
		tokens["uk"] = string(sub[1])
		return

	}

}

func (pcs *BaiduPCS) GenerateRequestQuery(mode, link string, params map[string]string) (res map[string]string) {
	res = make(map[string]string)
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"Referer":    "https://pan.baidu.com/disk/home",
	}
	if mode == "POST" {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, mode, link, params, headers)
	if panError != nil {
		res["ErrMsg"] = "未知错误"
		return
	}
	defer dataReadCloser.Close()
	body, err := ioutil.ReadAll(dataReadCloser)

	if err != nil {
		res["ErrMsg"] = "未知错误"
		return
	}
	if !gjson.Valid(string(body)) {
		res["ErrMsg"] = "返回值json解析错误"
		return
	}
	errno := gjson.Get(string(body), `errno`).Int()
	if errno != 0 {
		res["ErrMsg"] = "获取分享项元数据错误"
		if mode == "POST" && errno == 12 {
			path := gjson.Get(string(body), `info.0.path`).String()
			_, file := filepath.Split(path)
			res["ErrMsg"] = fmt.Sprintf("当前目录下已有%s同名文件/文件夹", file)
		}
		return
	}
	if mode != "POST" {
		res["title"] = gjson.Get(string(body), `title`).String()
		res["filename"] = gjson.Get(string(body), `list.0.server_filename`).String()
		var fids_str string = "["
		fsids := gjson.Get(string(body), `list.#.fs_id`).Array()
		if len(fsids) > 1 {
			res["filename"] += "等多个文件"
		}
		for _, sid := range fsids {
			fids_str += sid.String() + ","
		}
		res["fs_id"] = fids_str[:len(fids_str)-1] + "]"
	}
	res["ErrMsg"] = fmt.Sprintf("%d", errno)
	return
}
