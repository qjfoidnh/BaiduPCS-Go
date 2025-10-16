package baidupcs

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/pcserror"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/pcstime"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

type (
	// OrderBy 排序字段
	OrderBy string
	// Order 升序降序
	Order string
)

const (
	// OrderByName 根据文件名排序
	OrderByName OrderBy = "name"
	// OrderByTime 根据时间排序
	OrderByTime OrderBy = "time"
	// OrderBySize 根据大小排序, 注意目录无大小
	OrderBySize OrderBy = "size"
	// OrderAsc 升序
	OrderAsc Order = "asc"
	// OrderDesc 降序
	OrderDesc Order = "desc"
)

type (
	// HandleFileDirectoryFunc 处理文件或目录的元信息, 返回值控制是否退出递归
	HandleFileDirectoryFunc func(depth int, fdPath string, fd *FileDirectory, pcsError pcserror.Error) bool

	// FileDirectory 文件或目录的元信息
	FileDirectory struct {
		FsID     int64  // fs_id
		AppID    int64  // app_id
		Path     string // 路径
		Filename string // 文件名 或 目录名
		Ctime    int64  // 创建日期
		Mtime    int64  // 修改日期
		MD5      string // md5 值
		BlockListJSON
		Size        int64  // 文件大小 (目录为0)
		Isdir       bool   // 是否为目录
		Ifhassubdir bool   // 是否含有子目录 (只对目录有效)
		PreBase     string // 真正的base目录

		Parent   *FileDirectory    // 父目录信息
		Children FileDirectoryList // 子目录信息
	}

	// FileDirectoryList FileDirectory 的 指针数组
	FileDirectoryList []*FileDirectory

	// fdJSON 用于解析远程JSON数据
	fdJSON struct {
		FsID     int64  `json:"fs_id"` // fs_id
		AppID    int64  `json:"app_id"`
		Path     string `json:"path"`            // 路径
		Filename string `json:"server_filename"` // 文件名 或 目录名
		Ctime    int64  `json:"ctime"`           // 创建日期
		Mtime    int64  `json:"mtime"`           // 修改日期
		MD5      string `json:"md5"`             // md5 值
		BlockListJSON
		Size           int64 `json:"size"` // 文件大小 (目录为0)
		IsdirInt       int8  `json:"isdir"`
		IfhassubdirInt int8  `json:"ifhassubdir"`

		// 对齐
		_ *fdJSON
		_ []*fdJSON
	}

	fdData struct {
		*pcserror.PCSErrInfo
		List FileDirectoryList
	}

	fdDataJSONExport struct {
		*pcserror.PCSErrInfo
		List []*fdJSON `json:"list"`
	}

	// OrderOptions 列文件/目录可选项
	OrderOptions struct {
		By    OrderBy
		Order Order
	}
)

var (
	// DefaultOrderOptions 默认的排序
	DefaultOrderOptions = &OrderOptions{
		By:    OrderByName,
		Order: OrderAsc,
	}

	defaultOrderOptionsStr = fmt.Sprint(DefaultOrderOptions)
)

// FilesDirectoriesMeta 获取单个文件/目录的元信息
func (pcs *BaiduPCS) FilesDirectoriesMeta(path string) (data *FileDirectory, pcsError pcserror.Error) {
	if path == "" {
		path = PathSeparator
	}

	fds, err := pcs.FilesDirectoriesBatchMeta(path)
	if err != nil {
		return nil, err
	}

	// 返回了多条元信息
	if len(fds) != 1 {
		return nil, &pcserror.PCSErrInfo{
			Operation: OperationFilesDirectoriesMeta,
			ErrType:   pcserror.ErrTypeOthers,
			Err:       errors.New("未知返回数据"),
		}
	}
	return fds[0], nil
}

// FilesDirectoriesBatchMeta 获取多个文件/目录的元信息
func (pcs *BaiduPCS) FilesDirectoriesBatchMeta(paths ...string) (data FileDirectoryList, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareFilesDirectoriesBatchMeta(paths...)
	if pcsError != nil {
		return nil, pcsError
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPCSErrorInfo(OperationFilesDirectoriesMeta)
	// 服务器返回数据进行处理
	jsonData := fdData{
		PCSErrInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationFilesDirectoriesMeta, dataReadCloser, (*fdDataJSONExport)(unsafe.Pointer(&jsonData)))
	if pcsError != nil {
		return
	}

	// 修复MD5
	jsonData.List.fixMD5()

	data = jsonData.List
	return
}

// FilesDirectoriesList 获取目录下的文件和目录列表
func (pcs *BaiduPCS) FilesDirectoriesList(path string, options *OrderOptions) (data FileDirectoryList, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareFilesDirectoriesList(path, options)
	if pcsError != nil {
		return nil, pcsError
	}

	defer dataReadCloser.Close()

	jsonData := fdData{
		PCSErrInfo: pcserror.NewPCSErrorInfo(OperationFilesDirectoriesList),
	}

	pcsError = pcserror.HandleJSONParse(OperationFilesDirectoriesList, dataReadCloser, (*fdDataJSONExport)(unsafe.Pointer(&jsonData)))
	if pcsError != nil {
		return nil, pcsError
	}

	// 修复MD5
	jsonData.List.fixMD5()

	data = jsonData.List
	return
}

// Search 按文件名搜索文件, 不支持查找目录
func (pcs *BaiduPCS) Search(targetPath, keyword string, recursive bool) (fdl FileDirectoryList, pcsError pcserror.Error) {
	if targetPath == "" {
		targetPath = PathSeparator
	}

	dataReadCloser, pcsError := pcs.PrepareSearch(targetPath, keyword, recursive)
	if pcsError != nil {
		return nil, pcsError
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPCSErrorInfo(OperationSearch)
	jsonData := fdData{
		PCSErrInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationSearch, dataReadCloser, (*fdDataJSONExport)(unsafe.Pointer(&jsonData)))
	if pcsError != nil {
		return
	}

	// 修复MD5
	jsonData.List.fixMD5()

	fdl = jsonData.List
	return
}

func (pcs *BaiduPCS) recurseList(path string, depth int, options *OrderOptions, prebase string, handleFileDirectoryFunc HandleFileDirectoryFunc) (fdl FileDirectoryList, ok bool) {
	fdl, pcsError := pcs.FilesDirectoriesList(path, options)
	if pcsError != nil {
		ok := handleFileDirectoryFunc(depth, path, nil, pcsError) // 传递错误
		return nil, ok
	}

	for k := range fdl {
		fdl[k].PreBase = prebase
		ok = handleFileDirectoryFunc(depth+1, fdl[k].Path, fdl[k], nil)
		if !ok {
			return
		}

		if !fdl[k].Isdir {
			continue
		}

		fdl[k].Children, ok = pcs.recurseList(fdl[k].Path, depth+1, options, filepath.Join(prebase, filepath.Base(fdl[k].Path)), handleFileDirectoryFunc)
		if !ok {
			return
		}
	}

	return fdl, true
}

// FilesDirectoriesRecurseList 递归获取目录下的文件和目录列表
func (pcs *BaiduPCS) FilesDirectoriesRecurseList(path string, options *OrderOptions, handleFileDirectoryFunc HandleFileDirectoryFunc) (data FileDirectoryList) {
	fd, pcsError := pcs.FilesDirectoriesMeta(path)
	if pcsError != nil {
		handleFileDirectoryFunc(0, path, nil, pcsError) // 传递错误
		return nil
	}

	if !fd.Isdir { // 不是一个目录
		handleFileDirectoryFunc(0, path, fd, nil)
		return FileDirectoryList{fd}
	} else {
		handleFileDirectoryFunc(0, path, fd, nil)
	}

	data, _ = pcs.recurseList(path, 0, options, filepath.Base(path), handleFileDirectoryFunc)
	return data
}

// fixMD5 尝试修复MD5字段
// 服务器返回的MD5字段不一定正确了, 即是BlockList只有一个md5
// MD5字段使用BlockList中的md5
func (f *FileDirectory) fixMD5() {
	if len(f.BlockList) != 1 {
		return
	}
	f.MD5 = f.BlockList[0]
}

func (f *FileDirectory) String() string {
	builder := &strings.Builder{}
	tb := pcstable.NewTable(builder)
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})

	if f.Isdir {
		tb.AppendBulk([][]string{
			[]string{"类型", "目录"},
			[]string{"目录路径", f.Path},
			[]string{"目录名称", f.Filename},
		})
	} else {
		var md5info string
		if len(f.BlockList) > 1 {
			md5info = "md5 (可能不正确)"
		} else {
			md5info = "md5 (截图请打码)"
		}
		tb.AppendBulk([][]string{
			[]string{"类型", "文件"},
			[]string{"文件路径", f.Path},
			[]string{"文件名称", f.Filename},
			[]string{"文件大小", strconv.FormatInt(f.Size, 10) + ", " + converter.ConvertFileSize(f.Size)},
			[]string{md5info, f.MD5},
		})
	}

	tb.Append([]string{"app_id", strconv.FormatInt(f.AppID, 10)})
	tb.Append([]string{"fs_id", strconv.FormatInt(f.FsID, 10)})
	tb.AppendBulk([][]string{
		[]string{"创建日期", pcstime.FormatTime(f.Ctime)},
		[]string{"修改日期", pcstime.FormatTime(f.Mtime)},
	})

	if f.Ifhassubdir {
		tb.Append([]string{"是否含有子目录", "true"})
	}

	tb.Render()
	return builder.String()
}

func (fl FileDirectoryList) fixMD5() {
	for _, v := range fl {
		v.fixMD5()
		v.MD5 = DecryptMD5(v.MD5)
	}
}

// TotalSize 获取目录下文件的总大小
func (fl FileDirectoryList) TotalSize() int64 {
	var size int64
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		size += fl[k].Size

		// 递归获取
		if fl[k].Children != nil {
			size += fl[k].Children.TotalSize()
		}
	}
	return size
}

// Count 获取文件总数和目录总数
func (fl FileDirectoryList) Count() (fileN, directoryN int64) {
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		if fl[k].Isdir {
			directoryN++
		} else {
			fileN++
		}

		// 递归获取
		if fl[k].Children != nil {
			fN, dN := fl[k].Children.Count()
			fileN += fN
			directoryN += dN
		}
	}
	return
}

// AllFilePaths 返回所有的网盘路径, 包括子目录
func (fl FileDirectoryList) AllFilePaths() (pcspaths []string) {
	fN, dN := fl.Count()
	pcspaths = make([]string, 0, fN+dN)
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		pcspaths = append(pcspaths, fl[k].Path)

		if fl[k].Children != nil {
			pcspaths = append(pcspaths, fl[k].Children.AllFilePaths()...)
		}
	}
	return
}
