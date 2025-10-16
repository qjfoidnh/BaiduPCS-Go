package baidupcs

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/netdisksign"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/multipartreader"
	"github.com/qjfoidnh/baidu-tools/tieba"
)

type (
	reqType int
)

const (
	reqTypePCS = iota
	reqTypePan
)

func handleRespClose(resp *http.Response) error {
	if resp != nil {
		return resp.Body.Close()
	}
	return nil
}

func handleRespStatusError(operation string, resp *http.Response) pcserror.Error {
	//errInfo := pcserror.NewPCSErrorInfo(operation)
	// http 响应错误处理
	switch resp.StatusCode / 100 {
	case 4, 5:
		errInfo := pcserror.DecodePCSJSONError(operation, resp.Body)
		resp.Body.Close()
		errInfo.SetRemoteError()
		//errInfo.SetNetError(fmt.Errorf("http 响应错误, %s", resp.Status))
		return errInfo
	}

	return nil
}

func (pcs *BaiduPCS) sendReqReturnResp(rt reqType, op, method, urlStr string, post interface{}, header map[string]string) (resp *http.Response, pcsError pcserror.Error) {
	if header == nil {
		header = map[string]string{}
	}

	var (
		_, uaok = header["User-Agent"]
	)

	if !uaok {
		switch rt {
		case reqTypePCS:
			header["User-Agent"] = pcs.pcsUA
		case reqTypePan:
			header["User-Agent"] = pcs.panUA
		}
	}

	resp, err := pcs.client.Req(method, urlStr, post, header)
	if err != nil {
		handleRespClose(resp)
		switch rt {
		case reqTypePCS:
			return nil, &pcserror.PCSErrInfo{
				Operation: op,
				ErrType:   pcserror.ErrTypeNetError,
				Err:       err,
			}
		case reqTypePan:
			return nil, &pcserror.PanErrorInfo{
				Operation: op,
				ErrType:   pcserror.ErrTypeNetError,
				Err:       err,
			}
		}
		panic("unreachable")
	}
	return resp, nil
}

func (pcs *BaiduPCS) sendReqReturnReadCloser(rt reqType, op, method, urlStr string, post interface{}, header map[string]string) (readCloser io.ReadCloser, pcsError pcserror.Error) {
	resp, pcsError := pcs.sendReqReturnResp(rt, op, method, urlStr, post, header)
	if pcsError != nil {
		return
	}
	return resp.Body, nil
}

// PrepareUK 获取用户 UK, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareUK() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()

	query := url.Values{}
	query.Set("need_selfinfo", "1")

	panURL := &url.URL{
		Scheme:   "https",
		Host:     PanBaiduCom,
		Path:     "api/user/getinfo",
		RawQuery: query.Encode(),
	}

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationGetUK, http.MethodGet, panURL.String(), nil, nil)
	return
}

// PreparePCSServers 获取推荐的pcs服务器URL
func (pcs *BaiduPCS) PreparePCSServers() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL := pcs.generatePCSURL("file", "locateupload", map[string]string{
		"upload_version": "2.0",
		"app_id":         PanAppID,
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationGetPCSServer, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationGetPCSServer, http.MethodGet, pcsURL.String(), nil, nil)
	return
}

// PrepareQuotaInfo 获取当前用户空间配额信息, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareQuotaInfo() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL := pcs.generatePCSURL("quota", "info")
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationQuotaInfo, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationQuotaInfo, http.MethodGet, pcsURL.String(), nil, nil)
	return
}

// PrepareFilesDirectoriesBatchMeta 获取多个文件/目录的元信息, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareFilesDirectoriesBatchMeta(paths ...string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	sendData, err := (&PathsListJSON{}).JSON(paths...)
	if err != nil {
		panic(OperationFilesDirectoriesMeta + ", json 数据构造失败, " + err.Error())
	}

	pcsURL := pcs.generatePCSURL("file", "meta")
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationFilesDirectoriesMeta, pcsURL)

	// 表单上传
	mr := multipartreader.NewMultipartReader()
	mr.AddFormField("param", bytes.NewReader(sendData))
	mr.CloseMultipart()

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationFilesDirectoriesMeta, http.MethodPost, pcsURL.String(), mr, nil)
	return
}

// PrepareFilesDirectoriesList 获取目录下的文件和目录列表, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareFilesDirectoriesList(path string, options *OrderOptions) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	if options == nil {
		options = DefaultOrderOptions
	}
	if path == "" {
		path = PathSeparator
	}

	pcsURL := pcs.generatePCSURL("file", "list", map[string]string{
		"path":  path,
		"by":    *(*string)(unsafe.Pointer(&options.By)),
		"order": *(*string)(unsafe.Pointer(&options.Order)),
		"limit": "0-2147483647",
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationFilesDirectoriesList, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationFilesDirectoriesList, http.MethodGet, pcsURL.String(), nil, nil)
	return
}

func (pcs *BaiduPCS) PrepareFilesDirectoriesDiff(cursor string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	//bdstoken, pcsError := pcs.BDSToken()
	//if pcsError != nil {
	//	return
	//}
	if cursor == "" {
		cursor = "null"
	}
	ns := netdisksign.NewLocateDownloadSign(pcs.uid, pcs.GetBDUSS())
	pcsURL := pcs.generatePanURL("batch/filediff", map[string]string{
		"cursor": cursor,
		//"bdstoken": bdstoken,
		"clienttype": "1",
	})
	paramsURL := ns.URLParam()
	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationGetCursorDiff, http.MethodGet, pcsURL.String()+"&"+paramsURL, nil, nil)
	return
}

func (pcs *BaiduPCS) PrepareBDStoken() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL := pcs.generatePanURL("gettemplatevariable", map[string]string{
		"clienttype": "0",
		"app_id":     string(pcs.appID),
		"fields":     `["bdstoken"]`,
	})
	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationGetBDSToken, http.MethodGet, pcsURL.String(), nil, nil)
	return
}

// PrepareSearch 按文件名搜索文件, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareSearch(targetPath, keyword string, recursive bool) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	var re string
	if recursive {
		re = "1"
	} else {
		re = "0"
	}
	pcsURL := pcs.generatePCSURL("file", "search", map[string]string{
		"path": targetPath,
		"wd":   keyword,
		"re":   re,
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationSearch, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationSearch, http.MethodGet, pcsURL.String(), nil, nil)
	return
}

// PrepareRemove 批量删除文件/目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRemove(paths ...string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	sendData, err := (&PathsListJSON{}).JSON(paths...)
	if err != nil {
		panic(OperationMove + ", json 数据构造失败, " + err.Error())
	}

	pcsURL := pcs.generatePCSURL("file", "delete")
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationRemove, pcsURL)

	// 表单上传
	mr := multipartreader.NewMultipartReader()
	mr.AddFormField("param", bytes.NewReader(sendData))
	mr.CloseMultipart()

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationRemove, http.MethodPost, pcsURL.String(), mr, nil)
	return
}

// PrepareMkdir 创建目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareMkdir(pcspath string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL := pcs.generatePCSURL("file", "mkdir", map[string]string{
		"path": pcspath,
	})
	baiduPCSVerbose.Infof("%s URL: %s", OperationMkdir, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationMkdir, http.MethodPost, pcsURL.String(), nil, nil)
	return
}

func (pcs *BaiduPCS) prepareCpMvOp(op string, cpmvJSON ...*CpMvJSON) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	var method string
	switch op {
	case OperationCopy:
		method = "copy"
	case OperationMove, OperationRename:
		method = "move"
	default:
		panic("Unknown operation: " + op)
	}

	sendData, err := (&CpMvListJSON{
		List: cpmvJSON,
	}).JSON()
	if err != nil {
		//json 数据生成失败
		panic(err)
	}

	pcsURL := pcs.generatePCSURL("file", method)
	baiduPCSVerbose.Infof("%s URL: %s\n", op, pcsURL)

	// 表单上传
	mr := multipartreader.NewMultipartReader()
	mr.AddFormField("param", bytes.NewReader(sendData))
	mr.CloseMultipart()

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, op, http.MethodPost, pcsURL.String(), mr, nil)
	return
}

// PrepareRename 重命名文件/目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRename(from, to string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	return pcs.prepareCpMvOp(OperationRename, &CpMvJSON{
		From: from,
		To:   to,
	})
}

// PrepareCopy 批量拷贝文件/目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareCopy(cpmvJSON ...*CpMvJSON) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	return pcs.prepareCpMvOp(OperationCopy, cpmvJSON...)
}

// PrepareMove 批量移动文件/目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareMove(cpmvJSON ...*CpMvJSON) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	return pcs.prepareCpMvOp(OperationMove, cpmvJSON...)
}

// prepareRapidUpload 秒传文件, 不进行文件夹检查
func (pcs *BaiduPCS) prepareRapidUpload(targetPath, contentMD5, sliceMD5, crc32 string, length int64) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	//bdstoken, pcsError := pcs.BDSToken()
	//if pcsError != nil {
	//	return
	//}
	pcsURL := pcs.generatePCSURL2("xpan/file", "create", map[string]string{
		"access_token": pcs.accessToken,
		//"bdstoken": bdstoken,
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationRapidUpload, pcsURL)
	post := map[string]string{
		"block_list": mergeStringList(contentMD5),
		"path":       targetPath,
		"size":       strconv.FormatInt(length, 10),
		"isdir":      "0",
		"rtype":      "3",
	}
	baiduPCSVerbose.Infof("%s URL: %s, Post: %v\n", OperationRapidUpload, pcsURL, post)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePan, OperationRapidUpload, http.MethodPost, pcsURL.String(), post, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// prepareRapidUploadV2 秒传文件接口2, 不进行文件夹检查
func (pcs *BaiduPCS) prepareRapidUploadV2(targetPath, uploadid, policy, contentMD5, sliceMD5, dataContent, crc32 string, offset, length, totalSize, dataTime int64, blockListMD5 []string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcsURL := pcs.generatePanURL("precreate", nil)
	post := map[string]string{
		"uploadid":     uploadid,
		"path":         targetPath,
		"target_path":  path.Dir(targetPath) + "/",
		"size":         strconv.FormatInt(totalSize, 10),
		"data_offset":  strconv.FormatInt(offset, 10),
		"isdir":        "0",
		"local_mtime":  strconv.FormatInt(dataTime, 10),
		"local_ctime":  strconv.FormatInt(dataTime, 10),
		"rtype":        policy,
		"checkexist":   "0",
		"autoinit":     "1",
		"content-md5":  contentMD5,
		"slice-md5":    sliceMD5,
		"data_time":    strconv.FormatInt(dataTime, 10),
		"data_length":  strconv.FormatInt(length, 10),
		"data_content": dataContent,
		"block_list":   mergeStringList(blockListMD5...),
		"mode":         "1",
	}
	baiduPCSVerbose.Infof("%s URL: %s, Post: %v\n", OperationRapidUpload, pcsURL, post)

	if uploadid == "" {
		delete(post, "uploadid")
	}

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePan, OperationRapidUpload, http.MethodPost, pcsURL.String(), post, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "*/*",
		"Connection":   "keep-alive",
	})
	return
}

func (pcs *BaiduPCS) prepareFakeRapidUploadV2(targetPath, policy string, dateTime int64, blockListMD5 []string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcsURL := pcs.generatePanURL("precreate", map[string]string{
		"app_id":  PanAppID,
		"channel": "1",
		"web":     "1",
	})
	post := map[string]string{
		"path":        targetPath,
		"target_path": path.Dir(targetPath) + "/",
		"local_mtime": strconv.FormatInt(dateTime, 10),
		"autoinit":    "1",
		"rtype":       policy,
		"block_list":  mergeStringList(blockListMD5...),
	}
	baiduPCSVerbose.Infof("%s URL: %s, Post: %v\n", OperationRapidUpload, pcsURL, post)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePan, OperationRapidUpload, http.MethodPost, pcsURL.String(), post, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "*/*",
		"Connection":   "keep-alive",
	})
	return
}

// PrepareRapidUpload 秒传文件旧接口, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRapidUpload(targetPath, contentMD5, sliceMD5, crc32 string, length int64) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsError = pcs.CheckIsdir(OperationRapidUpload, targetPath, "", length)
	if pcsError != nil {
		return nil, pcsError
	}

	return pcs.prepareRapidUpload(targetPath, contentMD5, sliceMD5, crc32, length)
}

// PrepareRapidUploadV2 秒传文件新接口, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRapidUploadV2(targetPath, policy, uploadid, contentMD5, sliceMD5, dataContent, crc32 string, offset, length, totalSize, dataTime int64, blockListMD5 []string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsError = pcs.CheckIsdir(OperationRapidUpload, targetPath, policy, totalSize)
	if pcsError != nil {
		return nil, pcsError
	}
	rtype := pcs.policyTortype(policy)
	return pcs.prepareRapidUploadV2(targetPath, uploadid, rtype, contentMD5, sliceMD5, dataContent, crc32, offset, length, totalSize, dataTime, blockListMD5)
}

func (pcs *BaiduPCS) PrepareFakeRapidUploadV2(targetPath, policy string, length, dataTime int64, blockListMD5 []string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsError = pcs.CheckIsdir(OperationRapidUpload, targetPath, policy, length)
	if pcsError != nil {
		return nil, pcsError
	}
	rtype := pcs.policyTortype(policy)
	return pcs.prepareFakeRapidUploadV2(targetPath, rtype, dataTime, blockListMD5)
}

// PrepareLocateDownload 获取下载链接, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareLocateDownload(pcspath string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	bduss := pcs.GetBDUSS()
	// 检测uid
	if pcs.uid == 0 {
		t, err := tieba.NewUserInfoByBDUSS(bduss)
		if err != nil {
			return nil, &pcserror.PCSErrInfo{
				Operation: OperationLocateDownload,
				ErrType:   pcserror.ErrTypeNetError,
				Err:       err,
			}
		}
		pcs.uid = t.Baidu.UID
	}

	ns := netdisksign.NewLocateDownloadSign(pcs.uid, bduss)
	pcsURL := &url.URL{
		Scheme: GetHTTPScheme(pcs.isHTTPS),
		Host:   pcs.URL().Host,
		Path:   "/rest/2.0/pcs/file",
		RawQuery: (url.Values{
			"ant":        []string{"1"},
			"check_blue": []string{"1"},
			"es":         []string{"1"},
			"esl":        []string{"1"},
			"app_id":     []string{PanAppID},
			"method":     []string{"locatedownload"},
			"path":       []string{pcspath},
			"ver":        []string{"4.0"},
			"clienttype": []string{"17"},
			"channel":    []string{"0"},
			"apn_id":     []string{"1_0"},
			"freeisp":    []string{"0"},
			"queryfree":  []string{"0"},
			"use":        []string{"0"},
		}).Encode() + "&" + ns.URLParam(),
	}
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationLocateDownload, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationLocateDownload, http.MethodPost, pcsURL.String(), nil, pcs.getPanUAHeader())
	return
}

// PrepareLocatePanAPIDownload 从百度网盘首页获取下载链接, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareLocatePanAPIDownload(fidList ...int64) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()
	// 初始化
	var (
		sign, err = pcs.ph.CacheSignature()
	)
	if err != nil {
		return nil, &pcserror.PanErrorInfo{
			Operation: OperationLocatePanAPIDownload,
			ErrType:   pcserror.ErrTypeOthers,
			Err:       err,
		}
	}

	panURL := pcs.generatePanURL("download", nil)
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationLocatePanAPIDownload, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationLocatePanAPIDownload, http.MethodPost, panURL.String(), map[string]string{
		"sign":      sign.Sign(),
		"timestamp": sign.Timestamp(),
		"fidlist":   mergeInt64List(fidList...),
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// PrepareUpload 上传单个文件, 只返回服务器响应数据和错误信息（分片上传中的预上传部分）
func (pcs *BaiduPCS) PrepareUpload(policy string, targetPath string, uploadFunc UploadFunc) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()

	pcsURL := pcs.generatePCSURL("file", "upload", map[string]string{
		"path":  targetPath,
		"ondup": strings.Replace(policy, "rsync", "overwrite", -1),
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationUpload, pcsURL)

	resp, err := uploadFunc(pcsURL.String(), pcs.client.Jar)
	if err != nil {
		handleRespClose(resp)
		return nil, &pcserror.PCSErrInfo{
			Operation: OperationUpload,
			ErrType:   pcserror.ErrTypeNetError,
			Err:       err,
		}
	}

	pcsError = handleRespStatusError(OperationUpload, resp)
	if pcsError != nil {
		return
	}

	return resp.Body, nil
}

// PrepareUploadTmpFile 分片上传—文件分片及上传, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareUploadTmpFile(uploadFunc UploadFunc) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL := pcs.generatePCSURL("file", "upload", map[string]string{
		"type": "tmpfile",
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationUploadTmpFile, pcsURL)

	resp, err := uploadFunc(pcsURL.String(), pcs.client.Jar)
	if err != nil {
		handleRespClose(resp)
		return nil, &pcserror.PCSErrInfo{
			Operation: OperationUploadTmpFile,
			ErrType:   pcserror.ErrTypeNetError,
			Err:       err,
		}
	}

	pcsError = handleRespStatusError(OperationUploadTmpFile, resp)
	if pcsError != nil {
		return
	}

	return resp.Body, nil
}

// PrepareUploadCreateSuperFile 分片上传—合并分片文件, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareUploadCreateSuperFile(uploadid, rtype string, fileSize int64, targetPath string, blockList []string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()

	panURL := pcs.generatePanURL("create", nil)

	baiduPCSVerbose.Infof("%s URL: %s\n", OperationUploadCreateSuperFile, panURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePan, OperationUploadCreateSuperFile, http.MethodPost, panURL.String(), map[string]string{
		"uploadid": uploadid,
		"path":     targetPath,
		"size":     strconv.FormatInt(fileSize, 10),
		"isdir":    "0",
		"rtype":    rtype,
		//"local_mtime": strconv.FormatInt(dataTime, 10),
		//"local_ctime": strconv.FormatInt(dataTime, 10),
		"block_list":  mergeStringList(blockList...),
		"target_path": path.Dir(targetPath),
		//"bdstoken":    bdstoken,
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// PrepareUploadPrecreate 分片上传—Precreate, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareUploadPrecreate(targetPath, contentMD5, sliceMD5, crc32 string, size int64, blockList []string) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()
	panURL := &url.URL{
		Scheme: "https",
		Host:   PanBaiduCom,
		Path:   "api/precreate",
	}
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationUploadPrecreate, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationUploadPrecreate, http.MethodPost, panURL.String(), map[string]string{
		"path":         targetPath,
		"size":         strconv.FormatInt(size, 10),
		"isdir":        "0",
		"block_list":   mergeStringList(blockList...),
		"autoinit":     "1",
		"content-md5":  contentMD5,
		"slice-md5":    sliceMD5,
		"contentCrc32": crc32,
		"rtype":        "2",
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// PrepareUploadSuperfile2 另一个上传接口
func (pcs *BaiduPCS) PrepareUploadSuperfile2(uploadid, targetPath string, partseq int, partOffset int64, uploadFunc UploadFunc) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL := pcs.generatePCSURL("superfile2", "upload", map[string]string{
		"type":       "tmpfile",
		"path":       targetPath,
		"partseq":    strconv.Itoa(partseq),
		"partoffset": strconv.FormatInt(partOffset, 10),
		"uploadid":   uploadid,
		"vip":        "1",
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationUploadSuperfile2, pcsURL)

	resp, err := uploadFunc(pcsURL.String(), pcs.client.Jar)
	if err != nil {
		handleRespClose(resp)
		return nil, &pcserror.PCSErrInfo{
			Operation: OperationUploadSuperfile2,
			ErrType:   pcserror.ErrTypeNetError,
			Err:       err,
		}
	}

	pcsError = handleRespStatusError(OperationUpload, resp)
	if pcsError != nil {
		return
	}
	return resp.Body, nil
}

// PrepareCloudDlAddTask 添加离线下载任务, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareCloudDlAddTask(sourceURL, savePath string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL2 := pcs.generatePCSURL2("services/cloud_dl", "add_task", map[string]string{
		"app_id":       PanAppID,
		"task_from":    "0",
		"selected_idx": "1",
		"save_path":    savePath,
		"source_url":   sourceURL,
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationCloudDlAddTask, pcsURL2)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationCloudDlAddTask, http.MethodPost, pcsURL2.String(), nil, nil)
	return
}

// PrepareCloudDlQueryTask 精确查询离线下载任务, 只返回服务器响应数据和错误信息,
// taskids 例子: 12123,234234,2344, 用逗号隔开多个 task_id
func (pcs *BaiduPCS) PrepareCloudDlQueryTask(taskIDs string) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL2 := pcs.generatePCSURL2("services/cloud_dl", "query_task", map[string]string{
		"app_id":   PanAppID,
		"op_type":  "1",
		"task_ids": taskIDs,
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationCloudDlQueryTask, pcsURL2)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationCloudDlQueryTask, http.MethodGet, pcsURL2.String(), nil, nil)
	return
}

// PrepareCloudDlListTask 查询离线下载任务列表, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareCloudDlListTask() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL2 := pcs.generatePCSURL2("services/cloud_dl", "list_task", map[string]string{
		"need_task_info": "1",
		"status":         "255",
		"start":          "0",
		"limit":          "1000",
		"app_id":         PanAppID,
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationCloudDlListTask, pcsURL2)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationCloudDlListTask, http.MethodPost, pcsURL2.String(), nil, nil)
	return
}

func (pcs *BaiduPCS) prepareCloudDlCDTask(operation, method string, taskID int64) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL2 := pcs.generatePCSURL2("services/cloud_dl", method, map[string]string{
		"app_id":  PanAppID,
		"task_id": strconv.FormatInt(taskID, 10),
	})
	baiduPCSVerbose.Infof("%s URL: %s\n", operation, pcsURL2)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, operation, http.MethodPost, pcsURL2.String(), nil, nil)
	return
}

// PrepareCloudDlCancelTask 取消离线下载任务, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareCloudDlCancelTask(taskID int64) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	return pcs.prepareCloudDlCDTask(OperationCloudDlCancelTask, "cancel_task", taskID)
}

// PrepareCloudDlDeleteTask 取消离线下载任务, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareCloudDlDeleteTask(taskID int64) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	return pcs.prepareCloudDlCDTask(OperationCloudDlDeleteTask, "delete_task", taskID)
}

// PrepareCloudDlClearTask 清空离线下载任务记录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareCloudDlClearTask() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()
	pcsURL2 := pcs.generatePCSURL2("services/cloud_dl", "clear_task")
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationCloudDlClearTask, pcsURL2)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationCloudDlClearTask, http.MethodPost, pcsURL2.String(), nil, nil)
	return
}

// PrepareSharePSet 私密分享文件, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareSharePSet(paths []string, pwd string, period int) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()
	panURL := &url.URL{
		Scheme: "https",
		Host:   PanBaiduCom,
		Path:   "share/pset",
	}
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationShareSet, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationShareSet, http.MethodPost, panURL.String(), map[string]string{
		"path_list":    mergeStringList(paths...),
		"schannel":     "4",
		"channel_list": "[]",
		"period":       strconv.Itoa(period),
		"pwd":          pwd,
		"share_type":   "9",
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// PrepareShareCancel 取消分享, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareShareCancel(shareIDs []int64) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()
	panURL := &url.URL{
		Scheme: "https",
		Host:   PanBaiduCom,
		Path:   "share/cancel",
	}

	baiduPCSVerbose.Infof("%s URL: %s\n", OperationShareCancel, panURL)

	ss := converter.SliceInt64ToString(shareIDs)
	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationShareCancel, http.MethodPost, panURL.String(), map[string]string{
		"shareid_list": "[" + strings.Join(ss, ",") + "]",
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// PrepareShareList 列出分享列表, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareShareList(page int) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()

	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("desc", "1")
	query.Set("order", "time")

	panURL := &url.URL{
		Scheme:   "https",
		Host:     PanBaiduCom,
		Path:     "share/record",
		RawQuery: query.Encode(),
	}
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationShareList, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationShareList, http.MethodGet, panURL.String(), nil, nil)
	return
}

// PrepareShareSURLInfo 获取分享的详细信息, 包含密码, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareShareSURLInfo(shareID int64) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()

	query := url.Values{}
	query.Set("shareid", strconv.FormatInt(shareID, 10))
	query.Set("sign", converter.ToString(netdisksign.ShareSURLInfoSign(shareID)))

	panURL := &url.URL{
		Scheme:   "https",
		Host:     PanBaiduCom,
		Path:     "share/surlinfoinrecord",
		RawQuery: query.Encode(),
	}
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationShareSURLInfo, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationShareSURLInfo, http.MethodGet, panURL.String(), nil, nil)
	return
}

// PrepareRecycleList 列出回收站文件列表, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRecycleList(page int) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()

	panURL := pcs.generatePanURL("recycle/list", map[string]string{
		"num":  "100",
		"page": strconv.Itoa(page),
	})

	baiduPCSVerbose.Infof("%s URL: %s\n", OperationRecycleList, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationRecycleList, http.MethodGet, panURL.String(), nil, nil)
	return
}

// PrepareRecycleRestore 还原回收站文件或目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRecycleRestore(fidList ...int64) (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()

	pcsURL := pcs.generatePCSURL("file", "restore")
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationRecycleRestore, pcsURL)

	fsIDList := make([]*FsIDJSON, 0, len(fidList))
	for k := range fidList {
		fsIDList = append(fsIDList, &FsIDJSON{
			FsID: fidList[k],
		})
	}
	fsIDListJSON := FsIDListJSON{
		List: fsIDList,
	}

	sendData, err := jsoniter.Marshal(&fsIDListJSON)
	if err != nil {
		panic(err)
	}

	// 表单上传
	mr := multipartreader.NewMultipartReader()
	mr.AddFormField("param", bytes.NewReader(sendData))
	mr.CloseMultipart()

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationRecycleRestore, http.MethodPost, pcsURL.String(), mr, nil)
	return
}

// PrepareRecycleDelete 删除回收站文件或目录, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRecycleDelete(fidList ...int64) (dataReadCloser io.ReadCloser, panError pcserror.Error) {
	pcs.lazyInit()

	panURL := pcs.generatePanURL("recycle/delete", nil)
	baiduPCSVerbose.Infof("%s URL: %s\n", OperationRecycleDelete, panURL)

	dataReadCloser, panError = pcs.sendReqReturnReadCloser(reqTypePan, OperationRecycleDelete, http.MethodPost, panURL.String(), map[string]string{
		"fidlist": mergeInt64List(fidList...),
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	return
}

// PrepareRecycleClear 清空回收站, 只返回服务器响应数据和错误信息
func (pcs *BaiduPCS) PrepareRecycleClear() (dataReadCloser io.ReadCloser, pcsError pcserror.Error) {
	pcs.lazyInit()

	pcsURL := pcs.generatePCSURL("file", "delete", map[string]string{
		"type": "recycle",
	})

	baiduPCSVerbose.Infof("%s URL: %s\n", OperationRecycleClear, pcsURL)

	dataReadCloser, pcsError = pcs.sendReqReturnReadCloser(reqTypePCS, OperationRecycleClear, http.MethodGet, pcsURL.String(), nil, nil)
	return
}
