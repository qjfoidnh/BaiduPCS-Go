package pcsapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

type (
	TransferLink struct {
		Url   string `json:"url,omitempty" form:"url.omitempty"`
		Token string `json:"token,omitempty" form:"token.omitempty"`
	}
	TransferStructre struct {
		Download bool           `json:"download,omitempty" form:"download,omitempty"`
		Collect  bool           `json:"collect,omitempty" form:"collect,omitempty"`
		Rname    bool           `json:"rname,omitempty" form:"rname,omitempty"`
		Links    []TransferLink `json:"links,omitempty" form:"links,omitempty"`
	}
)

// 解析并保存分享链接
func runTransfer(ctx *gin.Context) {
	// 设置默认值
	args := TransferStructre{}
	if err := ctx.ShouldBind(&args); err != nil {
		fmt.Printf("transfer command failed with error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	errs := []error{}
	resps := []map[string]string{}

	for _, link := range args.Links {
		l := link.Url
		if strings.Contains(l, "bdlink=") || !strings.Contains(l, "pan.baidu.com/") {
			// if err := RunRapidTransfer(l); err != nil {
			// 	errs = append(errs, err)
			// }
			errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, "秒传已不再被支持"))
			continue
		}
		if strings.Contains(l, "?pwd=") {
			link.Url = strings.Split(l, "?pwd=")[0]
			link.Token = strings.Split(l, "?pwd=")[1]
		}
		if link.Url[len(link.Url)-1:] == "/" {
			link.Url = link.Url[0 : len(link.Url)-1]
		}
		featurestrs := strings.Split(link.Url, "/")
		featurestr := featurestrs[len(featurestrs)-1]
		if strings.Contains(featurestr, "init?") {
			featurestr = "1" + strings.Split(featurestr, "=")[1]
		}
		if len(featurestr) > 23 || featurestr[0:1] != "1" || len(link.Token) != 4 {
			errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, "链接地址或提取码非法"))
			continue
		}
		pcs := pcscommand.GetBaiduPCS()
		tokens := pcs.AccessSharePage(featurestr, true)
		if tokens["ErrMsg"] != "0" {
			errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"]))
			continue
		}

		var verifyUrl string
		featuremap := make(map[string]string)
		featuremap["bdstoken"] = tokens["bdstoken"]
		featuremap["surl"] = featurestr[1:]
		if link.Token != "" {
			verifyUrl = pcs.GenerateShareQueryURL("verify", map[string]string{
				"shareid":    tokens["shareid"],
				"time":       strconv.Itoa(int(time.Now().UnixMilli())),
				"clienttype": "1",
				"uk":         tokens["share_uk"],
			}).String()
			res := pcs.PostShareQuery(verifyUrl, link.Url, map[string]string{
				"pwd":       link.Token,
				"vcode":     "null",
				"vcode_str": "null",
				"bdstoken":  tokens["bdstoken"],
			})
			if res["ErrMsg"] != "0" {
				errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"]))
				continue
			}
		}
		pcs.UpdatePCSCookies(true)

		tokens = pcs.AccessSharePage(featurestr, false)
		if tokens["ErrMsg"] != "0" {
			errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"]))
			tokens["error_message"] = "access share page failed"
			resps = append(resps, tokens)
			continue
		}

		featureMap := map[string]string{
			"bdstoken": tokens["bdstoken"],
			"root":     "1",
			"web":      "5",
			"app_id":   baidupcs.PanAppID,
			"shorturl": featurestr[1:],
			"channel":  "chunlei",
		}
		queryShareInfoUrl := pcs.GenerateShareQueryURL("list", featureMap).String()
		trans_metas := pcs.ExtractShareInfo(queryShareInfoUrl, tokens["shareid"], tokens["share_uk"], tokens["bdstoken"])

		if trans_metas["ErrMsg"] != "success" {
			errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, trans_metas["ErrMsg"]))
			trans_metas["error_message"] = "extract share url info failed"
			resps = append(resps, trans_metas)
			continue
		}
		trans_metas["path"] = pcscommand.GetActiveUser().Workdir
		if trans_metas["item_num"] != "1" && args.Collect {
			trans_metas["filename"] += "等文件"
			trans_metas["path"] = path.Join(pcscommand.GetActiveUser().Workdir, trans_metas["filename"])
			pcs.Mkdir(trans_metas["path"])
		}
		trans_metas["referer"] = "https://pan.baidu.com/s/" + featurestr
		pcs.UpdatePCSCookies(true)
		resp := pcs.GenerateRequestQuery("POST", trans_metas)
		if resp["ErrNo"] != "0" {
			errs = append(errs, fmt.Errorf("%s失败: %s", baidupcs.OperationShareFileSavetoLocal, resp["ErrMsg"]))
			resp["error_message"] = "query saved file failed"
			resps = append(resps, resp)
			//if resp["ErrNo"] == "4" {
			//	transMetas["shorturl"] = featureStr
			//	pcs.SuperTransfer(transMetas, resp["limit"]) // 试验性功能, 当前未启用
			//}
			return
		}
		if args.Collect {
			resp["filename"] = trans_metas["filename"]
		}
		fmt.Printf("%s成功, 保存了%s到当前目录\n", baidupcs.OperationShareFileSavetoLocal, resp["filename"])
		if args.Download {
			// 开始后台下载
			go func(link TransferLink) {
				fmt.Printf("分享链接:%s 已保存，即将开始下载\n", link.Url)
				paths := strings.Split(resp["filenames"], ",")
				paths = paths[0 : len(paths)-1]
				pcscommand.RunDownload(paths, nil)
			}(link)
		}
	}
	if len(errs) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"error":   errs,
			"results": resps,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"results": resps,
		})
	}
}

// RunRapidTransfer 执行秒传链接解析及保存
func RunRapidTransfer(link string) (err error) {
	if strings.Contains(link, "bdlink=") || strings.Contains(link, "bdpan://") {
		r, _ := regexp.Compile(`(bdlink=|bdpan://)([^\s]+)`)
		link1 := r.FindStringSubmatch(link)[2]
		var decodeBytes []byte
		decodeBytes, err = base64.StdEncoding.DecodeString(link1)
		if err != nil {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
			return
		}
		link = string(decodeBytes)
	}
	link = strings.TrimSpace(link)
	substrs := strings.SplitN(link, "#", 4)
	if len(substrs) == 4 {
		md5, slicemd5 := substrs[0], substrs[1]
		length, _ := strconv.ParseInt(substrs[2], 10, 64)
		filename := path.Join(pcscommand.GetActiveUser().Workdir, substrs[3])
		err = RunRapidUpload(filename, md5, slicemd5, "", length)
		return
	}
	substrs = strings.Split(link, "|")
	if len(substrs) == 4 {
		md5, slicemd5 := substrs[2], substrs[3]
		length, _ := strconv.ParseInt(substrs[1], 10, 64)
		filename := path.Join(pcscommand.GetActiveUser().Workdir, substrs[0])
		err = RunRapidUpload(filename, md5, slicemd5, "", length)
		return
	}
	fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
	return
}

func initRunTransfer(group *gin.RouterGroup) {
	group.POST("transfer", runTransfer)
}
