package pcscommand

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RunShareTransfer 执行分享链接转存到网盘
func RunShareTransfer(params []string, opt *baidupcs.TransferOption) {
	var link = params[0]
	// 支持整段分享文本(如 "通过网盘分享的文件: 链接: https://... 提取码: xxxx ..."),
	// 先从中提取出规范的分享链接, 再交给 net/url 解析
	if shareLinkRe := regexp.MustCompile(`(https?://pan\.baidu\.com/s/[0-9A-Za-z_-]+)(?:\?pwd=([0-9A-Za-z]{4}))?`); shareLinkRe.MatchString(link) {
		if m := shareLinkRe.FindStringSubmatch(link); len(m) >= 2 {
			link = m[1]
			if len(params) == 1 && len(m) == 3 && m[2] != "" {
				params = append(params, m[2])
			}
		}
	}
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接格式非法")
		return
	}
	queryParams := parsedURL.Query()
	extraCode := queryParams.Get("pwd")
	if len(params) == 1 {
		if strings.Contains(link, "bdlink=") || !strings.Contains(link, "pan.baidu.com/") {
			//RunRapidTransfer(link, opt.Rname)
			fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "秒传已不再被支持")
			return
		}
	} else if len(params) == 2 {
		extraCode = params[1]
	}

	featureStr := path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
	if featureStr == "init" {
		featureStr = "1" + queryParams.Get("surl")
	}
	if len(featureStr) > 23 || featureStr[0:1] != "1" || len(extraCode) != 4 {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接地址或提取码非法")
		return
	}
	pcs := GetBaiduPCS()
	tokens := pcs.AccessSharePage(featureStr, true)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}

	verifyUrl := pcs.GenerateShareQueryURL("verify", map[string]string{
		"shareid":    tokens["shareid"],
		"time":       strconv.Itoa(int(time.Now().UnixMilli())),
		"clienttype": "1",
		"uk":         tokens["share_uk"],
	}).String()
	res := pcs.PostShareQuery(verifyUrl, link, map[string]string{
		"pwd":       extraCode,
		"vcode":     "null",
		"vcode_str": "null",
		"bdstoken":  tokens["bdstoken"],
	})
	if res["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"])
		return
	}

	pcs.UpdatePCSCookies(true)

	tokens = pcs.AccessSharePage(featureStr, false)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}
	featureMap := map[string]string{
		"bdstoken": tokens["bdstoken"],
		"root":     "1",
		"web":      "5",
		"app_id":   baidupcs.PanAppID,
		"shorturl": featureStr[1:],
		"channel":  "chunlei",
	}
	queryShareInfoUrl := pcs.GenerateShareQueryURL("list", featureMap).String()
	transMetas := pcs.ExtractShareInfo(queryShareInfoUrl, tokens["shareid"], tokens["share_uk"], tokens["bdstoken"])

	if transMetas["ErrMsg"] != "success" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, transMetas["ErrMsg"])
		return
	}

	// 如果指定了 fs_id, 覆盖解析结果
	if opt.FsID != "" {
		transMetas["fs_id"] = "[" + opt.FsID + "]"
		transMetas["filename"] = "user_specified_file" // 名字可能不准确，但不影响转存
		transMetas["item_num"] = "1"
	}

	transMetas["path"] = GetActiveUser().Workdir
	if transMetas["item_num"] != "1" && opt.Collect {
		transMetas["filename"] += "等文件"
		transMetas["path"] = path.Join(GetActiveUser().Workdir, transMetas["filename"])
		pcs.Mkdir(transMetas["path"])
	}
	transMetas["referer"] = "https://pan.baidu.com/s/" + featureStr
	pcs.UpdatePCSCookies(true)
	resp := pcs.GenerateRequestQuery("POST", transMetas)
	if resp["ErrNo"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, resp["ErrMsg"])
		//if resp["ErrNo"] == "4" {
		//	transMetas["shorturl"] = featureStr
		//	pcs.SuperTransfer(transMetas, resp["limit"]) // 试验性功能, 当前未启用
		//}
		return
	}
	if opt.Collect {
		resp["filename"] = transMetas["filename"]
	}
	fmt.Printf("%s成功, 保存了%s到当前目录\n", baidupcs.OperationShareFileSavetoLocal, resp["filename"])
	if opt.Download {
		fmt.Println("10s后开始下载")
		time.Sleep(10 * time.Second)
		paths := strings.Split(resp["filenames"], ",")
		RunDownload(paths, nil)
	}
}
