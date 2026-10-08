package pcscommand

// transfer --list：只列分享文件清单、不转存。
//
// 流程与 transfer 前半段完全相同（访问分享页 → 校验提取码 → share/list root=1），
// 之后不 POST /share/transfer，而是对子目录按 uk+shareid+dir 递归列出，输出一行
// `SHARE-LIST-JSON {...}`。不改当前工作目录、不写网盘、不保存配置。
// 失败时输出与 transfer 相同的「分享链接转存到网盘失败: <原因>」，上层沿用同一套错误分类。

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/tidwall/gjson"
)

const (
	shareListJSONPrefix = "SHARE-LIST-JSON "
	shareListPageSize   = 100
	shareListMaxReqs    = 200   // 单次清单请求上限（含翻页），防超大分享刷接口
	shareListMaxEntries = 20000 // 条目上限
	shareListInterval   = 300 * time.Millisecond
)

type shareListEntry struct {
	FsID        int64  `json:"fs_id"`
	Path        string `json:"path"` // 相对分享根，以 / 开头
	Filename    string `json:"server_filename"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"isdir"`
	MD5         string `json:"md5,omitempty"`
	ServerMtime int64  `json:"server_mtime"`
}

type shareListResult struct {
	ShortURL  string           `json:"shorturl"`
	ShareID   string           `json:"shareid"`
	ShareUK   string           `json:"share_uk"`
	FileCount int              `json:"file_count"`
	DirCount  int              `json:"dir_count"`
	TotalSize int64            `json:"total_size"`
	Requests  int              `json:"requests"`
	Truncated bool             `json:"truncated"`
	SubdirErr int              `json:"subdir_errors"` // 列失败而跳过的子目录数
	Complete  bool             `json:"complete"`      // 未截断且子目录全部列出：清单可作完整性依据
	Files     []shareListEntry `json:"files"`
}

type shareLister struct {
	pcs      *baidupcs.BaiduPCS
	tokens   map[string]string
	referer  string
	requests int
	res      shareListResult
	rootDir  string // 分享根在分享者网盘中的父目录，用于算相对路径
}

func (l *shareLister) fetch(params map[string]string) ([]gjson.Result, string) {
	if l.requests >= shareListMaxReqs || len(l.res.Files) >= shareListMaxEntries {
		l.res.Truncated = true
		return nil, ""
	}
	if l.requests > 0 {
		time.Sleep(shareListInterval)
	}
	l.requests++
	u := l.pcs.GenerateShareQueryURL("list", params).String()
	body, errMsg := l.pcs.GetShareListBody(u, l.referer)
	if os.Getenv("PCS_SHARE_LIST_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[share-list debug] params=%v err=%q body=%.300s\n", debugParams(params), errMsg, body)
	}
	if errMsg != "" {
		return nil, errMsg
	}
	return gjson.GetBytes(body, "list").Array(), ""
}

func (l *shareLister) add(item gjson.Result) {
	p := item.Get("path").String()
	rel := p
	if l.rootDir != "" && l.rootDir != "/" && strings.HasPrefix(p, l.rootDir+"/") {
		rel = strings.TrimPrefix(p, l.rootDir)
	}
	e := shareListEntry{
		FsID:        item.Get("fs_id").Int(),
		Path:        rel,
		Filename:    item.Get("server_filename").String(),
		Size:        item.Get("size").Int(),
		IsDir:       item.Get("isdir").Int() == 1,
		MD5:         item.Get("md5").String(),
		ServerMtime: item.Get("server_mtime").Int(),
	}
	l.res.Files = append(l.res.Files, e)
	if e.IsDir {
		l.res.DirCount++
	} else {
		l.res.FileCount++
		l.res.TotalSize += e.Size
	}
}

// listRoot 翻页列顶层；baseParams 为 transfer 同款查询参数（shorturl/root=1/bdstoken…）。
func (l *shareLister) listRoot(baseParams map[string]string) ([]gjson.Result, string) {
	var all []gjson.Result
	for page := 1; ; page++ {
		params := map[string]string{}
		for k, v := range baseParams {
			params[k] = v
		}
		params["page"] = strconv.Itoa(page)
		params["num"] = strconv.Itoa(shareListPageSize)
		items, errMsg := l.fetch(params)
		if errMsg != "" {
			return nil, errMsg
		}
		all = append(all, items...)
		if len(items) < shareListPageSize || l.res.Truncated {
			return all, ""
		}
	}
}

// listDir 递归列子目录（uk+shareid+dir 形式，web 分享页同款）。
func (l *shareLister) listDir(dir string) string {
	for page := 1; ; page++ {
		items, errMsg := l.fetch(map[string]string{
			"uk":         l.tokens["share_uk"],
			"shareid":    l.tokens["shareid"],
			"dir":        dir,
			"page":       strconv.Itoa(page),
			"num":        strconv.Itoa(shareListPageSize),
			"order":      "other",
			"desc":       "1",
			"showempty":  "0",
			"web":        "1",
			"channel":    "chunlei",
			"clienttype": "0",
			"app_id":     baidupcs.PanAppID,
			"bdstoken":   l.tokens["bdstoken"],
		})
		if errMsg != "" {
			return errMsg
		}
		for _, item := range items {
			l.add(item)
			if item.Get("isdir").Int() == 1 {
				l.listDirBestEffort(item.Get("path").String())
			}
		}
		if len(items) < shareListPageSize || l.res.Truncated {
			return ""
		}
	}
}

// listDirBestEffort 子目录列失败不致命：记数并跳过，最终清单标 complete=false。
// 顶层已列出、链接确认有效，单个子目录接口抖动不应让整个清单作废。
func (l *shareLister) listDirBestEffort(dir string) {
	if errMsg := l.listDir(dir); errMsg != "" {
		l.res.SubdirErr++
		fmt.Fprintf(os.Stderr, "列子目录失败（跳过）: %s: %s\n", dir, errMsg)
	}
}

// runShareFileList 在 RunShareTransfer 完成分享页访问与提取码校验后调用。
func runShareFileList(pcs *baidupcs.BaiduPCS, featureStr string, tokens map[string]string, rootParams map[string]string) {
	l := &shareLister{
		pcs:     pcs,
		tokens:  tokens,
		referer: "https://pan.baidu.com/s/" + featureStr,
		res: shareListResult{
			ShortURL: featureStr[1:],
			ShareID:  tokens["shareid"],
			ShareUK:  tokens["share_uk"],
			Files:    []shareListEntry{},
		},
	}
	roots, errMsg := l.listRoot(rootParams)
	if errMsg != "" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, errMsg)
		return
	}
	if len(roots) > 0 {
		l.rootDir = path.Dir(roots[0].Get("path").String())
	}
	for _, item := range roots {
		l.add(item)
	}
	for _, item := range roots {
		if item.Get("isdir").Int() == 1 {
			l.listDirBestEffort(item.Get("path").String())
		}
	}
	l.res.Requests = l.requests
	l.res.Complete = !l.res.Truncated && l.res.SubdirErr == 0
	out, err := json.Marshal(l.res)
	if err != nil {
		fmt.Printf("%s失败: 清单序列化错误 %s\n", baidupcs.OperationShareFileSavetoLocal, err)
		return
	}
	fmt.Println(shareListJSONPrefix + string(out))
}

// debugParams 调试输出时隐去 bdstoken。
func debugParams(params map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range params {
		if k == "bdstoken" {
			v = "***"
		}
		out[k] = v
	}
	return out
}
