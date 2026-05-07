package baidupcs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"github.com/tidwall/gjson"
)

type (
	// ShareOption 分享可选项
	TransferOption struct {
		Download  bool   // 是否直接开始下载
		Collect   bool   // 多文件整合
		Rname     bool   // 随机改文件名
		NoCollect bool   // 不创建汇总目录
		BatchSize int    // 每批转存文件数
		Parallel  int    // 并发数
		SkipRange string // 跳过转存的目录范围，如 "0001-1000"
	}

	// ShareFileInfo 分享文件信息
	ShareFileInfo struct {
		FsID      int64  // 文件ID
		Filename  string // 文件名
		Path      string // 路径
		FileCount int    // 叶子目录中文件数
	}
)

func (pcs *BaiduPCS) GenerateShareQueryURL(subPath string, params map[string]string) *url.URL {
	shareURL := &url.URL{
		Scheme: GetHTTPScheme(true),
		Host:   PanBaiduCom,
		Path:   "/share/" + subPath,
	}
	uv := shareURL.Query()
	for key, value := range params {
		uv.Set(key, value)
	}

	shareURL.RawQuery = uv.Encode()
	return shareURL
}

func (pcs *BaiduPCS) GetShareFileList(shareID int64, shareUK, shortURL, bdstoken string, opt *TransferOption, recursive bool) (chan *ShareFileInfo, string) {
	return pcs.GetShareFileListEx(shareID, shareUK, shortURL, bdstoken, opt, recursive, false)
}

func (pcs *BaiduPCS) GetShareFileListEx(shareID int64, shareUK, shortURL, bdstoken string, opt *TransferOption, recursive, collectLeafDirs bool) (chan *ShareFileInfo, string) {
	visitedDirs := make(map[int64]bool)
	fileChan := make(chan *ShareFileInfo, 1000) // buffer to avoid blocking
	errChan := make(chan string, 1)
	go func() {
		defer close(fileChan)
		err, _ := pcs.getShareFileListRecursiveEx(shareID, shareUK, shortURL, bdstoken, recursive,
			collectLeafDirs, "/", 0, visitedDirs, 0,
			fileChan, opt)
		if err != "" {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return nil, err
	default:
		return fileChan, ""
	}
}

func (pcs *BaiduPCS) getShareFileListRecursive(shareID int64, shareUK, shortURL, bdstoken string, recursive bool, parentPath string, visitedDirs map[int64]bool, depth int, fileChan chan *ShareFileInfo, opt *TransferOption) string {
	if strings.Contains(parentPath, opt.SkipRange) {
		baiduPCSVerbose.Infof("跳过目录1: %s\n", parentPath)

		return "skipped"
	}
	err, _ := pcs.getShareFileListRecursiveEx(shareID, shareUK, shortURL, bdstoken, recursive,
		false, parentPath, 0, visitedDirs, depth,
		fileChan, opt)
	return err
}

func (pcs *BaiduPCS) getShareFileListRecursiveEx(shareID int64, shareUK, shortURL, bdstoken string, recursive,
	collectLeafDirs bool, parentPath string, currentFsID int64, visitedDirs map[int64]bool, depth int,
	fileChan chan *ShareFileInfo, opt *TransferOption) (string, int) {
	page := 1
	pageSize := 100

	if depth > 100 {
		return "", 0
	}

	totalCount := 0
	hasChildDir := false
	dirFileCount := 0

	for {
		rootVal := "1"
		if parentPath != "/" {
			rootVal = "0"
		}
		featureMap := map[string]string{
			"bdstoken": bdstoken,
			"root":     rootVal,
			"web":      "5",
			"app_id":   PanAppID,
			"shorturl": shortURL,
			"channel":  "chunlei",
			"page":     strconv.Itoa(page),
			"num":      strconv.Itoa(pageSize),
		}
		queryShareInfoUrl := pcs.GenerateShareQueryURL("list", featureMap).String()

		postData := map[string]string{
			"dir": parentPath,
		}

		if strings.Contains(parentPath, opt.SkipRange) {
			baiduPCSVerbose.Infof("跳过目录3: %s\n", parentPath)
			continue
		}

		dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodPost, queryShareInfoUrl, postData, map[string]string{
			"User-Agent":   requester.UserAgent,
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		})
		if panError != nil {
			return "获取文件列表失败: " + panError.GetError().Error(), totalCount
		}

		body, err := ioutil.ReadAll(dataReadCloser)
		dataReadCloser.Close()
		if err != nil {
			return "读取文件列表失败", totalCount
		}

		errno := gjson.Get(string(body), `errno`).Int()
		if errno != 0 {
			if errno == 4 {
				time.Sleep(time.Second)
				continue
			}
			return fmt.Sprintf("获取文件列表失败, 错误码%d", errno), totalCount
		}

		list := gjson.Get(string(body), `list`).Array()
		if len(list) == 0 {
			break
		}

		for _, item := range list {
			fsID := item.Get(`fs_id`).Int()
			filename := item.Get(`server_filename`).String()
			filePath := item.Get(`path`).String()
			isDir := item.Get(`isdir`).Int() == 1

			if filename == "." || filename == ".." {
				continue
			}

			if strings.Contains(filePath, opt.SkipRange) {
				baiduPCSVerbose.Infof("跳过目录4: %s\n", filePath)

				continue
			}

			if isDir && recursive {
				if visitedDirs[fsID] {
					continue
				}

				visitedDirs[fsID] = true
				hasChildDir = true

				subPath := filePath + "/"
				if !strings.HasSuffix(subPath, "/") {
					subPath = subPath + "/"
				}

				errMsg, subCount := pcs.getShareFileListRecursiveEx(shareID, shareUK, shortURL, bdstoken, true,
					collectLeafDirs, subPath, fsID, visitedDirs, depth+1,
					fileChan, opt)
				if errMsg != "" {
					return errMsg, totalCount
				}

				totalCount += subCount
			} else if !isDir {
				if collectLeafDirs && recursive {
					dirFileCount++
				} else if !collectLeafDirs {
					fileChan <- &ShareFileInfo{
						FsID:     fsID,
						Filename: filename,
						Path:     filePath,
					}
					totalCount++
				}
			}
		}

		if len(list) < pageSize {
			break
		}
		page++
	}

	if collectLeafDirs && recursive && currentFsID != 0 && !hasChildDir {
		dirPath := strings.TrimSuffix(parentPath, "/")
		fileChan <- &ShareFileInfo{
			FsID:      currentFsID,
			Filename:  path.Base(dirPath),
			Path:      dirPath,
			FileCount: dirFileCount,
		}
		totalCount++
	}

	return "", totalCount
}

func (pcs *BaiduPCS) ExtractShareInfo(shareURL, shardID, shareUK, bdstoken string) (res map[string]string) {
	res = make(map[string]string)
	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodGet, shareURL, nil, map[string]string{
		"User-Agent":   requester.UserAgent,
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	})
	if panError != nil {
		res["ErrMsg"] = "提交分享项查询请求时发生错误"
		return
	}
	defer dataReadCloser.Close()
	body, _ := ioutil.ReadAll(dataReadCloser)
	errno := gjson.Get(string(body), `errno`).Int()
	if errno != 0 {
		res["ErrMsg"] = fmt.Sprintf("未知错误, 错误码%d", errno)
		if errno == 8001 {
			res["ErrMsg"] = "已触发验证, 请稍后再试"
		}
		return
	}
	res["filename"] = gjson.Get(string(body), `list.0.server_filename`).String()
	fsidList := gjson.Get(string(body), `list.#.fs_id`).Array()
	var fidsStr string = "["
	for _, sid := range fsidList {
		fidsStr += sid.String() + ","
	}

	res["shareid"] = shardID
	res["from"] = shareUK
	res["bdstoken"] = bdstoken
	shareUrl := &url.URL{
		Scheme: GetHTTPScheme(true),
		Host:   PanBaiduCom,
		Path:   "/share/transfer",
	}
	uv := shareUrl.Query()
	uv.Set("app_id", PanAppID)
	uv.Set("channel", "chunlei")
	uv.Set("clienttype", "0")
	uv.Set("web", "1")
	for key, value := range res {
		uv.Set(key, value)
	}
	res["item_num"] = strconv.Itoa(len(fsidList))
	res["ErrMsg"] = "success"
	res["fs_id"] = fidsStr[:len(fidsStr)-1] + "]"
	shareUrl.RawQuery = uv.Encode()
	res["shareUrl"] = shareUrl.String()
	return
}

func (pcs *BaiduPCS) PostShareQuery(url string, referer string, data map[string]string) (res map[string]string) {
	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodPost, url, data, map[string]string{
		"User-Agent":   requester.UserAgent,
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Referer":      referer,
	})
	res = make(map[string]string)
	if panError != nil {
		res["ErrMsg"] = "提交分享项查询请求时发生错误"
		return
	}
	defer dataReadCloser.Close()
	body, _ := ioutil.ReadAll(dataReadCloser)
	errno := gjson.Get(string(body), `errno`).Int()
	if errno != 0 {
		res["ErrMsg"] = fmt.Sprintf("未知错误, 错误码%d", errno)
		if errno == -9 {
			res["ErrMsg"] = "提取码错误"
		}
		return
	}
	res["randsk"] = gjson.Get(string(body), `randsk`).String()
	res["ErrMsg"] = "0"
	return
}

func (pcs *BaiduPCS) AccessSharePage(featurestr string, first bool) (tokens map[string]string) {
	tokens = make(map[string]string)
	tokens["ErrMsg"] = "0"
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
	headers["Referer"] = "https://pan.baidu.com/disk/home"
	if !first {
		headers["Referer"] = fmt.Sprintf("https://pan.baidu.com/share/init?surl=%s", featurestr[1:])
	}
	shareLink := fmt.Sprintf("https://pan.baidu.com/s/%s", featurestr)

	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodGet, shareLink, nil, headers)

	if panError != nil {
		tokens["ErrMsg"] = "访问分享页失败"
		return
	}
	defer dataReadCloser.Close()
	body, _ := ioutil.ReadAll(dataReadCloser)
	notFoundFlag := strings.Contains(string(body), "platform-non-found")
	errorPageTitle := strings.Contains(string(body), "error-404")
	if errorPageTitle {
		tokens["ErrMsg"] = "页面不存在"
		return
	}
	if notFoundFlag {
		tokens["ErrMsg"] = "分享链接已失效"
		return
	} else {
		re, _ := regexp.Compile(`(\{.+?loginstate.+?\})\);`)
		sub := re.FindSubmatch(body)
		if len(sub) < 2 {
			tokens["ErrMsg"] = "请确认登录参数中已经包含了网盘STOKEN"
			return
		}
		tokens["bdstoken"] = gjson.Get(string(sub[1]), `bdstoken`).String()
		tokens["uk"] = gjson.Get(string(sub[1]), `uk`).String()
		tokens["share_uk"] = gjson.Get(string(sub[1]), `share_uk`).String()
		tokens["shareid"] = gjson.Get(string(sub[1]), `shareid`).String()
		return
	}

}

func (pcs *BaiduPCS) GenerateRequestQuery(mode string, params map[string]string) (res map[string]string) {
	res = make(map[string]string)
	res["ErrNo"] = "0"
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"Referer":    params["referer"],
	}
	if mode == "POST" {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	postdata := make(map[string]string)
	postdata["fsidlist"] = params["fs_id"]
	postdata["path"] = params["path"]
	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, mode, params["shareUrl"], postdata, headers)
	if panError != nil {
		res["ErrNo"] = "1"
		res["ErrMsg"] = "网络错误"
		return
	}
	defer dataReadCloser.Close()
	body, err := ioutil.ReadAll(dataReadCloser)
	if err != nil {
		res["ErrNo"] = "-1"
		res["ErrMsg"] = "未知错误"
		return
	}
	if !gjson.Valid(string(body)) {
		res["ErrNo"] = "2"
		res["ErrMsg"] = "返回json解析错误"
		return
	}
	errno := gjson.Get(string(body), `errno`).Int()
	if errno != 0 {
		res["ErrNo"] = "3"
		res["ErrMsg"] = "获取分享项元数据错误:" + string(body)
		if mode == "POST" && errno == 12 {
			path := gjson.Get(string(body), `info.0.path`).String()
			_, file := filepath.Split(path) // Should be path.Split here, but never mind~
			_errno := gjson.Get(string(body), `info.0.errno`).Int()
			targetFileNums := gjson.Get(string(body), `target_file_nums`).Int()
			targetFileNumsLimit := gjson.Get(string(body), `target_file_nums_limit`).Int()
			if targetFileNums > targetFileNumsLimit {
				res["ErrNo"] = "4"
				res["ErrMsg"] = fmt.Sprintf("转存文件数%d超过当前用户上限, 当前用户单次最大转存数%d", targetFileNums, targetFileNumsLimit)
				res["limit"] = fmt.Sprintf("%d", targetFileNumsLimit)
			} else if _errno == -30 {
				res["ErrNo"] = "9"
				res["ErrMsg"] = fmt.Sprintf("当前目录下已有%s同名文件/文件夹", file)
			} else {
				res["ErrMsg"] = fmt.Sprintf("未知错误, 错误代码%d", _errno)
			}
		} else if mode == "POST" && errno == 4 {
			res["ErrMsg"] = fmt.Sprintf("文件重复")
		}
		return
	}

	_, res["filename"] = filepath.Split(gjson.Get(string(body), `info.0.path`).String())
	filenames := gjson.Get(string(body), `info.#.path`).Array()
	filenamesStr := ""
	for _, _path := range filenames {
		filenamesStr += "," + path.Base(_path.String())
	}
	if len(filenamesStr) < 1 {
		res["filenames"] = "default" + pcsutil.GenerateRandomString(5)
	} else {
		res["filenames"] = filenamesStr[1:]
	}
	if len(gjson.Get(string(body), `info.#.fsid`).Array()) > 1 {
		res["filename"] += "等多个文件/文件夹"
	}
	return
}

func (pcs *BaiduPCS) SuperTransfer(params map[string]string, limit string) {
	//headers := map[string]string{
	//	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
	//	"Referer":    params["referer"],
	//}
	//limit_num, _ := strconv.Atoi(limit)
	//fsidlist_str := params["fs_id"]
	//fsidlist := strings.Split(fsidlist_str[1:len(fsidlist_str)-1], ",")
	//listUrl := &url.URL{
	//	Scheme: GetHTTPScheme(true),
	//	Host:   PanBaiduCom,
	//	Path:   "/share/list",
	//}
	//uv := listUrl.Query()
	//uv.Set("app_id", PanAppID)
	//uv.Set("channel", "chunlei")
	//uv.Set("clienttype", "0")
	//uv.Set("web", "1")
	//uv.Set("page", "1")
	//uv.Set("num", "100")
	//uv.Set("shorturl", params["shorturl"])
	//uv.Set("root", "1")
	//dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodGet, listUrl.String(), nil, headers)
	//if panError != nil {
	//	res["ErrNo"] = "1"
	//	res["ErrMsg"] = "网络错误"
	//	return
	//}
	//defer dataReadCloser.Close()
	//body, err := ioutil.ReadAll(dataReadCloser)
	//res["ErrNo"] = "-1"
	//if err != nil {
	//	res["ErrMsg"] = "未知错误"
	//	return
	//}
	return

}

func (pcs *BaiduPCS) BatchTransferShortURL(shortURL, referer, shareID, shareUK, bdstoken string, fsIDs []int64, targetPath string, filename string) (successCount int, errMsg string) {
	var fidsStr string = "["
	for _, sid := range fsIDs {
		fidsStr += strconv.FormatInt(sid, 10) + ","
	}
	if len(fsIDs) > 0 {
		fidsStr = fidsStr[:len(fidsStr)-1]
	}
	fidsStr += "]"

	shareUrl := pcs.GenerateShareQueryURL("transfer", map[string]string{
		"app_id":     PanAppID,
		"channel":    "chunlei",
		"clienttype": "0",
		"web":        "1",
		"bdstoken":   bdstoken,
		"shareid":    shareID,
		"from":       shareUK,
		"filename":   filename,
	})

	postdata := map[string]string{
		"fsidlist": fidsStr,
		"path":     targetPath,
	}

	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      referer,
	}

	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodPost, shareUrl.String(), postdata, headers)
	if panError != nil {
		return 0, "网络错误: " + panError.GetError().Error()
	}
	defer dataReadCloser.Close()

	body, err := ioutil.ReadAll(dataReadCloser)
	if err != nil {
		return 0, "读取响应失败"
	}

	if !gjson.Valid(string(body)) {
		return 0, "响应JSON解析错误"
	}

	errno := gjson.Get(string(body), `errno`).Int()
	if errno != 0 {
		errMsgStr := gjson.Get(string(body), `errmsg`).String()
		if errMsgStr == "" {
			errMsgStr = fmt.Sprintf("未知错误, 错误码%d", errno)
		}
		if errno == 12 {
			targetFileNums := gjson.Get(string(body), `target_file_nums`).Int()
			targetFileNumsLimit := gjson.Get(string(body), `target_file_nums_limit`).Int()
			if targetFileNums > 0 && targetFileNumsLimit > 0 {
				return 0, fmt.Sprintf("转存文件数%d超过当前用户上限, 当前用户单次最大转存数%d", targetFileNums, targetFileNumsLimit)
			}
		}
		return 0, errMsgStr
	}

	infoList := gjson.Get(string(body), `info`).Array()
	return len(infoList), ""
}
