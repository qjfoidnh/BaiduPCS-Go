package pcscommand

import (
	"encoding/base64"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RunShareTransfer 执行分享链接转存到网盘
func RunShareTransfer(params []string, opt *baidupcs.TransferOption) {
	var link string
	var extraCode string
	if len(params) == 1 {
		link = params[0]
		if strings.Contains(link, "bdlink=") || !strings.Contains(link, "pan.baidu.com/") {
			//RunRapidTransfer(link, opt.Rname)
			fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "秒传已不再被支持")
			return
		}
		extraCode = "none"
		if strings.Contains(link, "?pwd=") {
			extraCode = strings.Split(link, "?pwd=")[1]
			link = strings.Split(link, "?pwd=")[0]
		}
	} else if len(params) == 2 {
		link = params[0]
		extraCode = params[1]
	}
	if link[len(link)-1:] == "/" {
		link = link[0 : len(link)-1]
	}
	featureStrs := strings.Split(link, "/")
	featureStr := featureStrs[len(featureStrs)-1]
	if strings.Contains(featureStr, "init?") {
		featureStr = "1" + strings.Split(featureStr, "=")[1]
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

	if extraCode != "none" {
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
		fmt.Println("即将开始下载")
		paths := strings.Split(resp["filenames"], ",")
		paths = paths[0 : len(paths)-1]
		RunDownload(paths, nil)
	}
}

// RunRapidTransfer 执行秒传链接解析及保存
func RunRapidTransfer(link string, rnameOpt ...bool) {
	if strings.Contains(link, "bdlink=") || strings.Contains(link, "bdpan://") {
		r, _ := regexp.Compile(`(bdlink=|bdpan://)([^\s]+)`)
		link1 := r.FindStringSubmatch(link)[2]
		decodeBytes, err := base64.StdEncoding.DecodeString(link1)
		if err != nil {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
			return
		}
		link = string(decodeBytes)
	}
	rname := false
	if len(rnameOpt) > 0 {
		rname = rnameOpt[0]
	}
	link = strings.TrimSpace(link)
	substrs := strings.SplitN(link, "#", 4)
	if len(substrs) == 4 {
		md5, slicemd5 := substrs[0], substrs[1]
		size, _ := strconv.ParseInt(substrs[2], 10, 64)
		filename := path.Join(GetActiveUser().Workdir, randReplaceStr(substrs[3], rname))
		RunRapidUpload(filename, md5, slicemd5, size)
	} else if len(substrs) == 3 {
		md5 := substrs[0]
		size, _ := strconv.ParseInt(substrs[1], 10, 64)
		filename := path.Join(GetActiveUser().Workdir, randReplaceStr(substrs[2], rname))
		RunRapidUpload(filename, md5, "", size)
	} else {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
	}
	return
}
