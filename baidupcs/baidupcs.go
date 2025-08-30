// Package baidupcs BaiduPCS RESTful API 工具包
package baidupcs

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/expires/cachemap"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/internal/panhome"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsverbose"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
)

const (
	// OperationGetUK 获取UK
	OperationGetUK = "获取UK"
	// OperationGetBDSToken
	OperationGetBDSToken = "获取bdstoken"
	// OperationGetCursorDiff
	OperationGetCursorDiff = "获取cursor后文件列表diff信息"
	// OperationGetPCSServer 获取可用的psc服务器
	OperationGetPCSServer = "获取可用的PCS服务器列表"
	// OperationQuotaInfo 获取当前用户空间配额信息
	OperationQuotaInfo = "获取当前用户空间配额信息"
	// OperationFilesDirectoriesMeta 获取文件/目录的元信息
	OperationFilesDirectoriesMeta = "获取文件/目录的元信息"
	// OperationFilesDirectoriesList 获取目录下的文件列表
	OperationFilesDirectoriesList = "获取目录下的文件列表"
	// OperationSearch 搜索
	OperationSearch = "搜索"
	// OperationRemove 删除文件/目录
	OperationRemove = "删除文件/目录"
	// OperationMkdir 创建目录
	OperationMkdir = "创建目录"
	// OperationRename 重命名文件/目录
	OperationRename = "重命名文件/目录"
	// OperationCopy 拷贝文件/目录
	OperationCopy = "拷贝文件/目录"
	// OperationMove 移动文件/目录
	OperationMove = "移动文件/目录"
	// OperationRapidUpload 秒传文件
	OperationRapidUpload = "秒传文件"
	// OperationUpload 上传单个文件
	OperationUpload = "上传单个文件"
	// OperationUploadTmpFile 分片上传—文件分片及上传
	OperationUploadTmpFile = "分片上传—文件分片及上传"
	// OperationUploadCreateSuperFile 分片上传—合并分片文件
	OperationUploadCreateSuperFile = "分片上传—合并分片文件"
	// OperationUploadPrecreate 分片上传—Precreate
	OperationUploadPrecreate = "分片上传—Precreate"
	// OperationUploadSuperfile2 分片上传—Superfile2
	OperationUploadSuperfile2 = "分片上传—Superfile2"
	// OperationDownloadFile 下载单个文件
	OperationDownloadFile = "下载单个文件"
	// OperationDownloadStreamFile 下载流式文件
	OperationDownloadStreamFile = "下载流式文件"
	// OperationLocateDownload 获取下载链接
	OperationLocateDownload = "获取下载链接"
	// OperationLocatePanAPIDownload 从百度网盘首页获取下载链接
	OperationLocatePanAPIDownload = "获取下载链接2"
	// OperationCloudDlAddTask 添加离线下载任务
	OperationCloudDlAddTask = "添加离线下载任务"
	// OperationCloudDlQueryTask 精确查询离线下载任务
	OperationCloudDlQueryTask = "精确查询离线下载任务"
	// OperationCloudDlListTask 查询离线下载任务列表
	OperationCloudDlListTask = "查询离线下载任务列表"
	// OperationCloudDlCancelTask 取消离线下载任务
	OperationCloudDlCancelTask = "取消离线下载任务"
	// OperationCloudDlDeleteTask 删除离线下载任务
	OperationCloudDlDeleteTask = "删除离线下载任务"
	// OperationCloudDlClearTask 清空离线下载任务记录
	OperationCloudDlClearTask = "清空离线下载任务记录"
	// OperationShareSet 创建分享链接
	OperationShareSet = "创建分享链接"
	// OperationShareCancel 取消分享
	OperationShareCancel = "取消分享"
	// OperationShareList 列出分享列表
	OperationShareList = "列出分享列表"
	// OperationShareSURLInfo 获取分享详细信息
	OperationShareSURLInfo = "获取分享详细信息"
	// OperationShareFileSavetoLocal 分享链接转存到网盘
	OperationShareFileSavetoLocal = "分享链接转存到网盘"
	// OperationRapidLinkSavetoLocal 秒传链接转存到网盘
	OperationRapidLinkSavetoLocal = "秒传链接转存到网盘"
	// OperationRecycleList 列出回收站文件列表
	OperationRecycleList = "列出回收站文件列表"
	// OperationRecycleRestore 还原回收站文件或目录
	OperationRecycleRestore = "还原回收站文件或目录"
	// OperationRecycleDelete 删除回收站文件或目录
	OperationRecycleDelete = "删除回收站文件或目录"
	// OperationRecycleClear 清空回收站
	OperationRecycleClear = "清空回收站"

	// OperationExportFileInfo 导出文件信息
	OperationExportFileInfo = "导出文件信息"
	// OperationGetRapidUploadInfo 获取文件秒传信息
	OperationGetRapidUploadInfo = "获取文件秒传信息"
	// OperationFixMD5 修复文件md5
	OperationFixMD5 = "修复文件md5"
	// OperationMatchPathByShellPattern 通配符匹配文件路径
	OperationMatchPathByShellPattern = "通配符匹配文件路径"

	// PCSBaiduCom pcs api地址
	PCSBaiduCom = "pcs.baidu.com"
	// PanBaiduCom 网盘首页api地址
	PanBaiduCom = "pan.baidu.com"
	// YunBaiduCom 网盘首页api地址2
	YunBaiduCom = "yun.baidu.com"
	// PanAppID 百度网盘appid
	PanAppID = "250528"
	// NetdiskUA 网盘客户端ua
	NetdiskUA = "netdisk;P2SP;3.0.0.8;netdisk;11.12.3;ANG-AN00;android-android;10.0;JSbridge4.4.0;jointBridge;1.1.0;"
	// DotBaiduCom .baidu.com
	DotBaiduCom = ".baidu.com"
	// PathSeparator 路径分隔符
	PathSeparator = "/"
)

var (
	baiduPCSVerbose = pcsverbose.New("BAIDUPCS")

	baiduComURL = &url.URL{
		Scheme: "http",
		Host:   "pan.baidu.com",
	}

	baiduPcsComURL = &url.URL{
		Scheme: "http",
		Host:   "baidupcs.com",
	}
)

type (
	// BaiduPCS 百度 PCS API 详情
	BaiduPCS struct {
		appID       int                   // app_id
		isHTTPS     bool                  // 是否启用https
		uid         uint64                // 百度uid
		client      *requester.HTTPClient // http 客户端
		accessToken string                // accessToken
		pcsUA       string
		pcsAddr     string
		panUA       string
		isSetPanUA  bool
		fixPCSAddr  bool
		ph          *panhome.PanHome
		cacheOpMap  cachemap.CacheOpMap
	}

	userInfoJSON struct {
		*pcserror.PanErrorInfo
		Records []struct {
			Uk int64 `json:"uk"`
		} `json:"records"`
	}

	userVarJSON struct {
		*pcserror.PanErrorInfo
		Result struct {
			BDSToken string `json:"bdstoken"`
		} `json:"result"`
	}
)

// UpdatePCSCookies 去除重名Cookies, 同名cookeis只保留最新的
func (pcs *BaiduPCS) UpdatePCSCookies(reverse bool) {
	cookies := pcs.GetClient().Jar.Cookies(baiduComURL)
	new_cookies := make([]*http.Cookie, 0)
	usedKeys := make(map[string]string)
	if reverse {
		for i := len(cookies) - 1; i >= 0; i-- {
			cookie := cookies[i]
			_, exists := usedKeys[cookie.Name]
			if !exists {
				new_cookies = append(new_cookies, &http.Cookie{
					Name:   cookie.Name,
					Value:  cookie.Value,
					Domain: DotBaiduCom,
				})
				usedKeys[cookie.Name] = ""
			}
		}
		for i, j := 0, len(new_cookies)-1; i < j; i, j = i+1, j-1 {
			new_cookies[i], new_cookies[j] = new_cookies[j], new_cookies[i] //逆序
		}
	} else {
		for i := 0; i < len(cookies); i++ {
			cookie := cookies[i]
			_, exists := usedKeys[cookie.Name]
			if !exists {
				new_cookies = append(new_cookies, &http.Cookie{
					Name:   cookie.Name,
					Value:  cookie.Value,
					Domain: DotBaiduCom,
				})
				usedKeys[cookie.Name] = ""
			}
		}
	}
	client := pcs.GetClient()
	client.ResetCookiejar()
	// fmt.Println(len(new_cookies))
	// for _, c := range new_cookies {
	// 	fmt.Println(c.Value)
	// }
	client.Jar.SetCookies(baiduComURL, new_cookies)

}

// NewPCS 提供app_id, 百度BDUSS, 返回 BaiduPCS 对象
func NewPCS(appID int, bduss string) *BaiduPCS {
	client := requester.NewHTTPClient()
	client.ResetCookiejar()
	client.Jar.SetCookies(baiduComURL, []*http.Cookie{
		&http.Cookie{
			Name:   "BDUSS",
			Value:  bduss,
			Domain: DotBaiduCom,
		},
	})

	return &BaiduPCS{
		appID:  appID,
		client: client,
	}
}

// NewPCSWithClient 提供app_id, 自定义客户端, 返回 BaiduPCS 对象
func NewPCSWithClient(appID int, client *requester.HTTPClient) *BaiduPCS {
	pcs := &BaiduPCS{
		appID:  appID,
		client: client,
	}
	return pcs
}

// NewPCSWithCookieStr 提供app_id, cookie 字符串, 返回 BaiduPCS 对象
func NewPCSWithCookieStr(appID int, cookieStr string) *BaiduPCS {
	pcs := &BaiduPCS{
		appID:  appID,
		client: requester.NewHTTPClient(),
	}

	cookies := requester.ParseCookieStr(cookieStr)
	for _, cookie := range cookies {
		cookie.Domain = DotBaiduCom
	}

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(baiduComURL, cookies)
	pcs.client.SetCookiejar(jar)

	return pcs
}

func (pcs *BaiduPCS) lazyInit() {
	if pcs.client == nil {
		pcs.client = requester.NewHTTPClient()
	}
	if pcs.ph == nil {
		pcs.ph = panhome.NewPanHome(pcs.client)
	}
	if !pcs.isSetPanUA {
		pcs.panUA = NetdiskUA
	}
}

// GetClient 获取当前的http client
func (pcs *BaiduPCS) GetClient() *requester.HTTPClient {
	pcs.lazyInit()
	return pcs.client
}

// GetBDUSS 获取BDUSS
func (pcs *BaiduPCS) GetBDUSS() (bduss string) {
	if pcs.client == nil || pcs.client.Jar == nil {
		return ""
	}
	cookies := pcs.client.Jar.Cookies(baiduComURL)
	for _, cookie := range cookies {
		if cookie.Name == "BDUSS" {
			return cookie.Value
		}
	}
	return ""
}

// GetBAIDUID 获取BAIDUID
func (pcs *BaiduPCS) GetBAIDUID() (baiduid string) {
	if pcs.client == nil || pcs.client.Jar == nil {
		return ""
	}
	cookies := pcs.client.Jar.Cookies(baiduComURL)
	for _, cookie := range cookies {
		if cookie.Name == "BAIDUID" {
			return cookie.Value
		}
	}
	return ""
}

func (pcs *BaiduPCS) GetPCSAddr() string {
	return pcs.pcsAddr
}

// SetAPPID 设置app_id
func (pcs *BaiduPCS) SetAPPID(appID int) {
	pcs.appID = appID
}

// SetUID 设置百度UID
// 只有locatedownload才需要设置此项
func (pcs *BaiduPCS) SetUID(uid uint64) {
	pcs.uid = uid
}

// SetaccessToken 设置秒传转存用的accesstoken
func (pcs *BaiduPCS) SetaccessToken(accessToken string) {
	pcs.accessToken = accessToken
}

// SetStoken 设置stoken
func (pcs *BaiduPCS) SetStoken(stoken string) {
	pcs.lazyInit()
	if pcs.client.Jar == nil {
		pcs.client.ResetCookiejar()
	}

	pcs.client.Jar.SetCookies(baiduComURL, []*http.Cookie{
		&http.Cookie{
			Name:   "STOKEN",
			Value:  stoken,
			Domain: DotBaiduCom,
		},
	})
}

// SetSboxtkn 设置sboxtkn
func (pcs *BaiduPCS) SetSboxtkn(sboxtkn string) {
	pcs.lazyInit()
	if pcs.client.Jar == nil {
		pcs.client.ResetCookiejar()
	}

	pcs.client.Jar.SetCookies(baiduComURL, []*http.Cookie{
		&http.Cookie{
			Name:   "SBOXTKN",
			Value:  sboxtkn,
			Domain: DotBaiduCom,
		},
	})
}

// SetPCSUserAgent 设置 PCS User-Agent
func (pcs *BaiduPCS) SetPCSUserAgent(ua string) {
	pcs.pcsUA = ua
}

// SetPCSAddr 设置 PCS 服务器地址
func (pcs *BaiduPCS) SetPCSAddr(addr string) {
	if addr != "" {
		pcs.pcsAddr = addr
	}
}

// SetPanUserAgent 设置 Pan User-Agent
func (pcs *BaiduPCS) SetPanUserAgent(ua string) {
	pcs.panUA = ua
	pcs.isSetPanUA = true
}

// SetHTTPS 是否启用https连接
func (pcs *BaiduPCS) SetHTTPS(https bool) {
	pcs.isHTTPS = https
}

func (pcs *BaiduPCS) SetStaticPCSAddr(static bool) {
	pcs.fixPCSAddr = static
}

// URL 返回 url
func (pcs *BaiduPCS) URL() *url.URL {
	host := pcs.pcsAddr
	if host == "" {
		host = PCSBaiduCom
	}
	return &url.URL{
		Scheme: GetHTTPScheme(pcs.isHTTPS),
		Host:   host,
	}
}

func (pcs *BaiduPCS) getPanUAHeader() (header map[string]string) {
	return map[string]string{
		"User-Agent": pcs.panUA,
	}
}

func (pcs *BaiduPCS) generatePCSURL(subPath, method string, param ...map[string]string) *url.URL {
	pcsURL := pcs.URL()
	pcsURL.Path = "/rest/2.0/pcs/" + subPath

	uv := pcsURL.Query()
	uv.Set("app_id", strconv.Itoa(pcs.appID))
	uv.Set("method", method)
	for k := range param {
		for k2 := range param[k] {
			uv.Set(k2, param[k][k2])
		}
	}

	pcsURL.RawQuery = uv.Encode()
	return pcsURL
}

func (pcs *BaiduPCS) generatePCSURL2(subPath, method string, param ...map[string]string) *url.URL {
	pcsURL2 := &url.URL{
		Scheme: GetHTTPScheme(pcs.isHTTPS),
		Host:   PanBaiduCom,
		Path:   "/rest/2.0/" + subPath,
	}

	uv := pcsURL2.Query()
	uv.Set("method", method)
	for k := range param {
		for k2 := range param[k] {
			uv.Set(k2, param[k][k2])
		}
	}

	pcsURL2.RawQuery = uv.Encode()
	return pcsURL2
}

func (pcs *BaiduPCS) generatePanURL(subPath string, param map[string]string) *url.URL {
	panURL := url.URL{
		Scheme: GetHTTPScheme(pcs.isHTTPS),
		Host:   PanBaiduCom,
		Path:   "/api/" + subPath,
	}

	if param != nil {
		uv := url.Values{}
		for k := range param {
			uv.Set(k, param[k])
		}
		panURL.RawQuery = uv.Encode()
	}
	return &panURL
}

// UK 获取用户 UK
func (pcs *BaiduPCS) UK() (uk int64, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUK()
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationGetUK)
	jsonData := userInfoJSON{
		PanErrorInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationGetUK, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	if len(jsonData.Records) != 1 {
		errInfo.ErrType = pcserror.ErrTypeOthers
		errInfo.Err = errors.New("Unknown remote data")
		return 0, errInfo
	}

	return jsonData.Records[0].Uk, nil
}

func (pcs *BaiduPCS) BDSToken() (bdstoken string, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareBDStoken()
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationGetBDSToken)
	jsonData := userVarJSON{
		PanErrorInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationGetBDSToken, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	if jsonData.Result.BDSToken == "" {
		errInfo.ErrType = pcserror.ErrTypeOthers
		errInfo.Err = errors.New("Unknown remote data")
		return "", errInfo
	}

	return jsonData.Result.BDSToken, nil
}
