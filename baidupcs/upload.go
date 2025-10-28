package baidupcs

import (
	"errors"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	// MaxUploadBlockSize 上传的文件分片最大大小
	MaxUploadBlockSize = 64 * converter.MB
	// MiddleUploadBlockSize 上传的文件分片中等大小
	MiddleUploadBlockSize = 16 * converter.MB
	// MinUploadBlockSize 上传的文件分片最小大小
	MinUploadBlockSize = 4 * converter.MB
	// RecommendedUploadSize 推荐的最高文件上传大小
	RecommendedUploadSize = 32 * converter.GB
	// MaxUploadSize 目前支持的最大文件大小
	MaxUploadSize = 128 * converter.GB
	// SliceMD5Size 计算 slice-md5 所需的长度
	SliceMD5Size = 256 * converter.KB
	// EmptyContentMD5 空串的md5
	EmptyContentMD5 = "d41d8cd98f00b204e9800998ecf8427e"
	// MiddleUploadThreshold 中等分片对应的文件大小
	MiddleUploadThreshold = 8 * converter.GB
	// MaxUploadThreshold 最大分片对应的文件大小
	MaxUploadThreshold = 32 * converter.GB
	// MinCheckLeftSpaceThreshold 需要检查剩余空间是否足够的最小文件大小
	MinCheckLeftSpaceThreshold = 64 * converter.MB
)

var (
	// ErrUploadMD5NotFound 未找到md5
	ErrUploadMD5NotFound = errors.New("unknown response data, md5 not found")
	// ErrUploadSavePathFound 未找到保存路径
	ErrUploadSavePathFound = errors.New("unknown response data, file saved path not found")
	// ErrUploadSeqNotMatch 服务器返回的上传队列不匹配
	ErrUploadSeqNotMatch = errors.New("服务器返回的上传队列不匹配")
	// ErrUploadMD5Unknown 服务器无匹配文件/秒传未生效
	ErrUploadMD5Unknown = errors.New("服务器无匹配文件/秒传未生效")
	// ErrUploadFileExists 文件或目录已存在
	ErrUploadFileExists = errors.New("文件已存在")

	// FakeBlockListMD5 虚假秒传时的BlockList
	fakeBlockListMD5 = []string{"5910a591dd8fc18c32a8f3df4fdc1761", "a5fc157d78e6ad1c7e114b056c92821e"}
)

type (
	// UploadFunc 上传文件处理函数
	UploadFunc func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error)

	// RapidUploadInfo 文件秒传信息
	RapidUploadInfo struct {
		Filename      string
		ContentLength int64
		ContentMD5    string
		SliceMD5      string
		ContentCrc32  string
	}

	uploadTmpFileJSON struct {
		MD5 string `json:"md5"`
		*pcserror.PCSErrInfo
	}

	uploadPrecreateJSON struct {
		ReturnType int    `json:"return_type"` // 1上传, 2秒传
		UploadID   string `json:"uploadid"`
		//BlockList  []int  `json:"block_list"`
		*pcserror.PanErrorInfo
		fdJSON `json:"info"`
	}

	// UploadSeq 分片上传顺序
	UploadSeq struct {
		Seq   int
		Block string
	}

	// PrecreateInfo 预提交文件消息返回数据
	PrecreateInfo struct {
		IsRapidUpload bool
		UploadID      string
		UploadSeqList []*UploadSeq
	}

	PCSServer struct {
		ServerAddr string `json:"server"`
	}

	PCSInfo struct {
		*pcserror.PCSErrInfo
		Host    string       `json:"host"`
		Server  []string     `json:"server"`
		Servers []*PCSServer `json:"servers"`
	}
)

// RapidUpload 秒传文件
func (pcs *BaiduPCS) RapidUpload(targetPath, policy, uploadid, contentMD5, sliceMD5, dataContent, crc32 string, offset, length, totalSize, dataTime int64, blockListMD5 []string) (pcsError pcserror.Error, jsonData uploadPrecreateJSON) {
	defer func() {
		if pcsError == nil {
			// 更新缓存
			pcs.deleteCache([]string{path.Dir(targetPath)})
		}
	}()
	pcsError, jsonData = pcs.rapidUploadV2(targetPath, policy, uploadid, strings.ToLower(contentMD5), strings.ToLower(sliceMD5), dataContent, crc32, offset, length, totalSize, dataTime, blockListMD5)
	return
}

// FakeRapidUpload 只precreate不进行秒传
func (pcs *BaiduPCS) FakeRapidUpload(targetPath, policy string, length int64) (pcsError pcserror.Error, jsonData uploadPrecreateJSON) {
	defer func() {
		if pcsError == nil {
			// 更新缓存
			pcs.deleteCache([]string{path.Dir(targetPath)})
		}
	}()
	if length <= MinUploadBlockSize {
		pcsError, jsonData = pcs.fakeRapidUploadV2(targetPath, policy, length, time.Now().Unix(), fakeBlockListMD5[0:1])
		return
	}
	pcsError, jsonData = pcs.fakeRapidUploadV2(targetPath, policy, length, time.Now().Unix(), fakeBlockListMD5)
	return
}

func (pcs *BaiduPCS) rapidUploadV2(targetPath, policy, uploadid, contentMD5, sliceMD5, dataContent, crc32 string, offset, length, totalSize, dataTime int64, blockListMD5 []string) (pcsError pcserror.Error, jsonData uploadPrecreateJSON) {
	dataReadCloser, pcsError := pcs.PrepareRapidUploadV2(targetPath, policy, uploadid, contentMD5, sliceMD5, dataContent, crc32, offset, length, totalSize, dataTime, blockListMD5)
	if pcsError != nil {
		return
	}
	defer dataReadCloser.Close()
	jsonData = uploadPrecreateJSON{
		PanErrorInfo: pcserror.NewPanErrorInfo(OperationRapidUpload),
	}
	pcsError = pcserror.HandleJSONParse(OperationUpload, dataReadCloser, &jsonData)
	return pcsError, jsonData
}

func (pcs *BaiduPCS) fakeRapidUploadV2(targetPath, policy string, length, dateTime int64, blockListMD5 []string) (pcsError pcserror.Error, jsonData uploadPrecreateJSON) {
	dataReadCloser, pcsError := pcs.PrepareFakeRapidUploadV2(targetPath, policy, length, dateTime, blockListMD5)
	if pcsError != nil {
		return
	}
	defer dataReadCloser.Close()
	jsonData = uploadPrecreateJSON{
		PanErrorInfo: pcserror.NewPanErrorInfo(OperationRapidUpload),
	}
	pcsError = pcserror.HandleJSONParse(OperationUpload, dataReadCloser, &jsonData)
	return pcsError, jsonData
}

// RapidUploadNoCheckDir 秒传文件, 不进行目录检查, 会覆盖掉同名的目录!
func (pcs *BaiduPCS) RapidUploadNoCheckDir(targetPath, contentMD5, sliceMD5, crc32 string, length int64) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.prepareRapidUpload(targetPath, contentMD5, sliceMD5, crc32, length)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	pcsError = pcserror.DecodePCSJSONError(OperationRapidUpload, dataReadCloser)
	if pcsError != nil {
		return
	}

	return nil
}

// UploadTmpFile 分片上传—文件分片及上传
func (pcs *BaiduPCS) UploadTmpFile(uploadid, targetPath string, partseq int, partOffset int64, uploadFunc UploadFunc) (md5 string, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUploadSuperfile2(uploadid, targetPath, partseq, partOffset, uploadFunc)
	if pcsError != nil {
		return "", pcsError
	}

	defer dataReadCloser.Close()

	// 数据处理
	jsonData := uploadTmpFileJSON{
		PCSErrInfo: pcserror.NewPCSErrorInfo(OperationUploadTmpFile),
	}

	pcsError = pcserror.HandleJSONParse(OperationUploadTmpFile, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	// 未找到md5
	if jsonData.MD5 == "" {
		jsonData.PCSErrInfo.ErrType = pcserror.ErrTypeInternalError
		jsonData.PCSErrInfo.Err = ErrUploadMD5NotFound
		return "", jsonData.PCSErrInfo
	}

	return jsonData.MD5, nil
}

// UploadCreateSuperFile 分片上传—合并分片文件
func (pcs *BaiduPCS) UploadCreateSuperFile(uploadid, policy string, fileSize int64, targetPath string, checksumMap map[int]string) (panError pcserror.Error) {
	blockList := sortBlockList(checksumMap)
	rtype := pcs.policyTortype(policy)
	dataReadCloser, pcsError := pcs.PrepareUploadCreateSuperFile(uploadid, rtype, fileSize, targetPath, blockList)
	if pcsError != nil {
		return pcsError
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.DecodePanJSONError(OperationUploadCreateSuperFile, dataReadCloser)
	if errInfo != nil {
		return errInfo
	}

	// 更新缓存, targetPath取了dir所以不受重命名策略影响
	pcs.deleteCache([]string{path.Dir(targetPath)})
	return nil
}

// GetRandomPCSHost 随机获取一个可用的pcs地址
func (pcs *BaiduPCS) GetRandomPCSHost() (pcsError pcserror.Error, pcsHost string) {
	if pcs.fixPCSAddr {
		return
	}
	dataReadCloser, pcsError := pcs.PreparePCSServers()
	if pcsError != nil {
		return
	}
	defer dataReadCloser.Close()
	pcsInfo := &PCSInfo{
		PCSErrInfo: pcserror.NewPCSErrorInfo(OperationGetPCSServer),
	}
	pcsError = pcserror.HandleJSONParse(OperationGetPCSServer, dataReadCloser, pcsInfo)
	if pcsError != nil {
		return
	}
	pcsHostList := make([]string, 0)
	if len(pcsInfo.Servers) > 0 {
		for _, server := range pcsInfo.Servers {
			if strings.Contains(server.ServerAddr, "-") {
				parsedURL, err := url.Parse(server.ServerAddr)
				if err != nil {
					continue
				}
				pcsHostList = append(pcsHostList, parsedURL.Hostname())
			}
		}
	} else if len(pcsInfo.Server) > 0 {
		pcsHostList = pcsInfo.Server
	}
	pcsHost = RandomElement(pcsHostList)
	return
}
