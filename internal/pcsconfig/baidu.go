package pcsconfig

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/baidu-tools/tieba"
)

var (
	//ErrNoSuchBaiduUser 未登录任何百度帐号
	ErrNoSuchBaiduUser = errors.New("no such baidu user")
	//ErrBaiduUserNotFound 未找到百度帐号
	ErrBaiduUserNotFound = errors.New("baidu user not found")
)

//BaiduBase Baidu基
type BaiduBase struct {
	UID  uint64 `json:"uid"`  // 百度ID对应的uid
	Name string `json:"name"` // 真实ID
}

// Baidu 百度帐号对象
type Baidu struct {
	BaiduBase
	Sex string  `json:"sex"` // 性别
	Age float64 `json:"age"` // 帐号年龄

	BDUSS   string `json:"bduss"`
	PTOKEN  string `json:"ptoken"`
	STOKEN  string `json:"stoken"`
	BAIDUID string `json:"baiduid"`
	SBOXTKN string `json:"sboxtkn"`
	COOKIES string `json:"cookies"`

	AccessToken string `json:"accesstoken"`

	Workdir string `json:"workdir"` // 工作目录
}

// BaiduPCS 初始化*baidupcs.BaiduPCS
func (baidu *Baidu) BaiduPCS() *baidupcs.BaiduPCS {
	pcs := baidupcs.NewPCS(Config.AppID, baidu.BDUSS)
	pcs.SetStoken(baidu.STOKEN)
	if baidu.SBOXTKN != "" {
		pcs.SetSboxtkn(baidu.SBOXTKN)
	}
	if strings.Contains(baidu.COOKIES, "STOKEN=") && baidu.STOKEN == "" {
		// 未显式指定stoken则从cookies中读取
		pcs = baidupcs.NewPCSWithCookieStr(Config.AppID, baidu.COOKIES)
	}
	pcs.SetHTTPS(Config.EnableHTTPS)
	pcs.SetPCSUserAgent(Config.PCSUA)
	pcs.SetPanUserAgent(Config.PanUA)
	pcs.SetUID(baidu.UID)
	pcs.SetaccessToken(baidu.AccessToken)
	return pcs
}

// GetSavePath 根据提供的网盘文件路径 pcspath, 返回本地储存路径,
// 返回绝对路径, 获取绝对路径出错时才返回相对路径...
func (baidu *Baidu) GetSavePath(pcspath string) string {
	dirStr := filepath.Join(Config.SaveDir, fmt.Sprintf("%d_%s", baidu.UID, converter.TrimPathInvalidChars(baidu.Name)), pcspath)
	dir, err := filepath.Abs(dirStr)
	if err != nil {
		dir = filepath.Clean(dirStr)
	}
	return dir
}

// PathJoin 合并工作目录和相对路径p, 若p为绝对路径则忽略
func (baidu *Baidu) PathJoin(p string) string {
	if path.IsAbs(p) {
		return p
	}
	return path.Join(baidu.Workdir, p)
}

// BaiduUserList 百度帐号列表
type BaiduUserList []*Baidu

// NewUserInfoByBDUSS 检测BDUSS有效性, 同时获取百度详细信息 (无法获取 ptoken 和 stoken)
func NewUserInfoByBDUSS(bduss string) (b *Baidu, err error) {
	t, err := tieba.NewUserInfoByBDUSS(bduss)
	if err != nil {
		return nil, err
	}

	b = &Baidu{
		BaiduBase: BaiduBase{
			UID:  t.Baidu.UID,
			Name: t.Baidu.Name,
		},
		Sex:     t.Baidu.Sex,
		Age:     t.Baidu.Age,
		BDUSS:   bduss,
		Workdir: "/",
	}
	return b, nil
}

// NewUserInfoByInput 不检测BDUSS有效性, 手动设置百度详细信息 (只适用tieba.baidu.com用户信息接口不可用的情况)
func NewUserInfoByInput(bduss, name string) (b *Baidu, err error) {
	b = &Baidu{
		BaiduBase: BaiduBase{
			UID:  24,
			Name: name,
		},
		Sex:     "default",
		Age:     24,
		BDUSS:   bduss,
		Workdir: "/",
	}
	return b, nil
}

// String 格式输出百度帐号列表
func (bl *BaiduUserList) String() string {
	builder := &strings.Builder{}

	tb := pcstable.NewTable(builder)
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	tb.SetHeader([]string{"#", "uid", "用户名", "性别", "age"})

	for k, baiduInfo := range *bl {
		tb.Append([]string{strconv.Itoa(k), strconv.FormatUint(baiduInfo.UID, 10), baiduInfo.Name, baiduInfo.Sex, fmt.Sprint(baiduInfo.Age)})
	}

	tb.Render()

	return builder.String()
}
