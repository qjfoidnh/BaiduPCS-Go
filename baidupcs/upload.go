package baidupcs

import (
	"errors"
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"net/http"
	"path"
)

const (
	// MaxUploadBlockSize 最大上传的文件分片大小
	MaxUploadBlockSize = 2 * converter.GB
	// MinUploadBlockSize 最小的上传的文件分片大小
	MinUploadBlockSize = 4 * converter.MB
	// MaxRapidUploadSize 秒传文件支持的最大文件大小
	MaxRapidUploadSize = 20 * converter.GB
	// RecommendUploadBlockSize 推荐的上传的文件分片大小
	RecommendUploadBlockSize = 1 * converter.GB
	// SliceMD5Size 计算 slice-md5 所需的长度
	SliceMD5Size = 256 * converter.KB
	// EmptyContentMD5 空串的md5
	EmptyContentMD5 = "d41d8cd98f00b204e9800998ecf8427e"
)

var (
	// ErrUploadMD5NotFound 未找到md5
	ErrUploadMD5NotFound = errors.New("unknown response data, md5 not found")
	// ErrUploadSavePathFound 未找到保存路径
	ErrUploadSavePathFound = errors.New("unknown response data, file saved path not found")
	// ErrUploadSeqNotMatch 服务器返回的上传队列不匹配
	ErrUploadSeqNotMatch = errors.New("服务器返回的上传队列不匹配")
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

	uploadJSON struct {
		*PathJSON
		*pcserror.PCSErrInfo
	}

	uploadTmpFileJSON struct {
		MD5 string `json:"md5"`
		*pcserror.PCSErrInfo
	}

	uploadPrecreateJSON struct {
		ReturnType int    `json:"return_type"` // 1上传, 2秒传
		UploadID   string `json:"uploadid"`
		BlockList  []int  `json:"block_list"`
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

	uploadSuperfile2JSON struct {
		MD5 string `json:"md5"`
		*pcserror.PCSErrInfo
	}
)

// RapidUpload 秒传文件
func (pcs *BaiduPCS) RapidUpload(targetPath, contentMD5, sliceMD5, crc32 string, length int64) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareRapidUpload(targetPath, contentMD5, sliceMD5, crc32, length)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	pcsError = pcserror.DecodePCSJSONError(OperationRapidUpload, dataReadCloser)
	if pcsError != nil {
		return
	}

	// 更新缓存
	pcs.deleteCache([]string{path.Dir(targetPath)})
	return nil
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

// Upload 上传单个文件
func (pcs *BaiduPCS) Upload(targetPath string, uploadFunc UploadFunc) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUpload(targetPath, uploadFunc)
	if pcsError != nil {
		return pcsError
	}

	defer dataReadCloser.Close()

	// 数据处理
	jsonData := uploadJSON{
		PCSErrInfo: pcserror.NewPCSErrorInfo(OperationUpload),
	}

	pcsError = pcserror.HandleJSONParse(OperationUpload, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	if jsonData.Path == "" {
		jsonData.PCSErrInfo.ErrType = pcserror.ErrTypeInternalError
		jsonData.PCSErrInfo.Err = ErrUploadSavePathFound
		return jsonData.PCSErrInfo
	}

	// 更新缓存
	pcs.deleteCache([]string{path.Dir(targetPath)})
	return nil
}

// UploadTmpFile 分片上传—文件分片及上传
func (pcs *BaiduPCS) UploadTmpFile(uploadFunc UploadFunc) (md5 string, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUploadTmpFile(uploadFunc)
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
func (pcs *BaiduPCS) UploadCreateSuperFile(checkDir bool, targetPath string, blockList ...string) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUploadCreateSuperFile(checkDir, targetPath, blockList...)
	if pcsError != nil {
		return pcsError
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.DecodePCSJSONError(OperationUploadCreateSuperFile, dataReadCloser)
	if errInfo != nil {
		return errInfo
	}

	// 更新缓存
	pcs.deleteCache([]string{path.Dir(targetPath)})
	return nil
}

// UploadPrecreate 分片上传—Precreate,
// 支持检验秒传
func (pcs *BaiduPCS) UploadPrecreate(targetPath, contentMD5, sliceMD5, crc32 string, size int64, bolckList ...string) (precreateInfo *PrecreateInfo, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUploadPrecreate(targetPath, contentMD5, sliceMD5, crc32, size, bolckList...)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationUploadPrecreate)
	jsonData := uploadPrecreateJSON{
		PanErrorInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationUploadPrecreate, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	switch jsonData.ReturnType {
	case 1: // 上传
		seqLen := len(jsonData.BlockList)
		if seqLen != len(bolckList) {
			errInfo.ErrType = pcserror.ErrTypeRemoteError
			errInfo.Err = ErrUploadSeqNotMatch
			return nil, errInfo
		}

		seqList := make([]*UploadSeq, 0, seqLen)
		for k, seq := range jsonData.BlockList {
			seqList = append(seqList, &UploadSeq{
				Seq:   seq,
				Block: bolckList[k],
			})
		}
		return &PrecreateInfo{
			UploadID:      jsonData.UploadID,
			UploadSeqList: seqList,
		}, nil

	case 2: // 秒传
		return &PrecreateInfo{
			IsRapidUpload: true,
		}, nil

	default:
		panic("unknown returntype")
	}
}

// UploadSuperfile2 分片上传—Superfile2
func (pcs *BaiduPCS) UploadSuperfile2(uploadid, targetPath string, partseq int, partOffset int64, uploadFunc UploadFunc) (md5sum string, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareUploadSuperfile2(uploadid, targetPath, partseq, partOffset, uploadFunc)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	jsonData := uploadSuperfile2JSON{
		PCSErrInfo: pcserror.NewPCSErrorInfo(OperationUploadSuperfile2),
	}

	pcsError = pcserror.HandleJSONParse(OperationUploadSuperfile2, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	return jsonData.MD5, nil
}
