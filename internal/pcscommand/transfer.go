package pcscommand

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/iikira/BaiduPCS-Go/baidupcs"
)

// RunShareTransfer 执行分享链接转存到网盘
func RunShareTransfer(params []string) {
	var link string
	var extracode string
	if len(params) == 1 {
		link = params[0]
		if strings.Contains(link, "bdlink=") || !strings.Contains(link, "pan.baidu.com/") {
			RunRapidTransfer(link)
			return
		}
		extracode = "none"
	} else if len(params) == 2 {
		link = params[0]
		extracode = params[1]
	}
	if link[len(link)-1:] == "/" {
		link = link[0 : len(link)-1]
	}
	featurestrs := strings.Split(link, "/")
	featurestr := featurestrs[len(featurestrs)-1]
	if strings.Contains(featurestr, "init?") {
		featurestr = "1" + strings.Split(featurestr, "=")[1]
	}
	if len(featurestr) != 23 || featurestr[0:1] != "1" || len(extracode) != 4 {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接地址或提取码非法")
		return
	}
	pcs := GetBaiduPCS()
	tokens := pcs.AccessSharePage(featurestr, true)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}
	// pcs.UpdatePCSCookies(true)
	var vefiryurl string
	featuremap := make(map[string]string)
	featuremap["surl"] = featurestr[1:]
	featuremap["bdstoken"] = tokens["bdstoken"]
	if extracode != "none" {

		vefiryurl = pcs.GenerateShareQueryURL("verify", featuremap).String()
		res := pcs.PostShareQuery(vefiryurl, featurestr[1:], map[string]string{
			"pwd":       extracode,
			"vcode":     "",
			"vcode_str": "",
		})
		if res["ErrMsg"] != "0" {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"])
			return
		}
	}
	pcs.UpdatePCSCookies(true)

	tokens = pcs.AccessSharePage(featurestr, false)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}
	metajsonstr := tokens["metajson"]
	trans_metas := pcs.ExtractShareInfo(metajsonstr)

	if trans_metas["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, trans_metas["ErrMsg"])
		return
	}
	trans_metas["path"] = GetActiveUser().Workdir
	trans_metas["referer"] = "https://pan.baidu.com/s/" + featurestr
	pcs.UpdatePCSCookies(true)
	resp := pcs.GenerateRequestQuery("POST", trans_metas)
	if resp["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, resp["ErrMsg"])
		return
	}
	fmt.Printf("%s成功, 保存了%s到当前目录\n", baidupcs.OperationShareFileSavetoLocal, resp["filename"])
}

// RunRapidTransfer 执行秒传链接解析及保存
func RunRapidTransfer(link string) {
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
	substrs := strings.Split(link, "#")
	if len(substrs) == 4 {
		md5 := strings.ToLower(substrs[0])
		slicemd5 := strings.ToLower(substrs[1])
		length, _ := strconv.ParseInt(substrs[2], 10, 64)
		filename := filepath.Join(GetActiveUser().Workdir, substrs[3])
		RunRapidUpload(filename, md5, slicemd5, "", length)
		return
	}
	substrs = strings.Split(link, "|")
	if len(substrs) == 4 {
		md5 := strings.ToLower(substrs[2])
		slicemd5 := strings.ToLower(substrs[3])
		length, _ := strconv.ParseInt(substrs[1], 10, 64)
		filename := filepath.Join(GetActiveUser().Workdir, substrs[0])
		RunRapidUpload(filename, md5, slicemd5, "", length)
		return
	}
	fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
}
