package pcsconfig

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"os"
	"strconv"
)

// ActiveUser 获取当前登录的用户
func (c *PCSConfig) ActiveUser() *Baidu {
	if c.activeUser == nil {
		return &Baidu{}
	}
	return c.activeUser
}

// ActiveUserBaiduPCS 获取当前登录的用户的baidupcs.BaiduPCS
func (c *PCSConfig) ActiveUserBaiduPCS() *baidupcs.BaiduPCS {
	if c.pcs == nil {
		c.pcs = c.ActiveUser().BaiduPCS()
	}
	return c.pcs
}

func (c *PCSConfig) httpClientWithUA(ua string) *requester.HTTPClient {
	client := requester.NewHTTPClient()
	client.SetHTTPSecure(c.EnableHTTPS)
	client.SetUserAgent(ua)
	return client
}

// HTTPClient 返回设置好的 HTTPClient
func (c *PCSConfig) HTTPClient() *requester.HTTPClient {
	return c.httpClientWithUA(c.UserAgent)
}

// PCSHTTPClient 返回设置好的 PCS HTTPClient
func (c *PCSConfig) PCSHTTPClient() *requester.HTTPClient {
	return c.httpClientWithUA(c.PCSUA)
}

// PanHTTPClient 返回设置好的 Pan HTTPClient
func (c *PCSConfig) PanHTTPClient() *requester.HTTPClient {
	return c.httpClientWithUA(c.PanUA)
}

// NumLogins 获取登录的用户数量
func (c *PCSConfig) NumLogins() int {
	return len(c.BaiduUserList)
}

// AverageParallel 返回平均的下载最大并发量
func (c *PCSConfig) AverageParallel() int {
	return AverageParallel(c.MaxParallel, c.MaxDownloadLoad)
}

// PrintTable 输出表格
func (c *PCSConfig) PrintTable() {
	tb := pcstable.NewTable(os.Stdout)
	tb.SetHeader([]string{"名称", "值", "建议值", "描述"})
	tb.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	tb.AppendBulk([][]string{
		[]string{"appid", fmt.Sprint(c.AppID), "", "百度 PCS 应用ID"},
		[]string{"cache_size", converter.ConvertFileSize(int64(c.CacheSize), 2), "1KB ~ 256KB", "下载缓存, 如果硬盘占用高或下载速度慢, 请尝试调大此值"},
		[]string{"max_parallel", strconv.Itoa(c.MaxParallel), "1 ~ 20", "下载总最大并发量, 非svip不可>1"},
		[]string{"max_upload_parallel", strconv.Itoa(c.MaxUploadParallel), "1 ~ 100", "上传单文件最大并发量"},
		[]string{"max_download_load", strconv.Itoa(c.MaxDownloadLoad), "1 ~ 5", "同时进行下载文件的最大数量"},
		[]string{"max_download_rate", showMaxRate(c.MaxDownloadRate), "", "限制最大下载速度, 0代表不限制"},
		[]string{"max_upload_rate", showMaxRate(c.MaxUploadRate), "", "限制最大上传速度, 0代表不限制"},
		[]string{"max_upload_load", strconv.Itoa(c.MaxUploadLoad), "1 ~ 4", "同时进行上传文件的最大数量"},
		[]string{"savedir", c.SaveDir, "", "下载文件的储存目录"},
		[]string{"enable_https", fmt.Sprint(c.EnableHTTPS), "true", "启用 https"},
		[]string{"force_login_username", fmt.Sprint(c.ForceLogin), "留空", "强制登录指定用户名, 适用于tieba用户信息接口不可用的情况, 如登录正常请留空"},
		[]string{"ignore_illegal", fmt.Sprint(c.IgnoreIllegal), "false", "关闭上传文件的文件名非法字符检查"},
		[]string{"upload_policy", fmt.Sprint(c.UPolicy), baidupcs.SkipPolicy, fmt.Sprintf("上传遇到重名文件时的处理策略, %s(默认，跳过)、%s(覆盖)、%s(仅跳过大小未变化的文件其余覆盖)",
			baidupcs.SkipPolicy, baidupcs.OverWritePolicy, baidupcs.RsyncPolicy)},
		[]string{"user_agent", c.UserAgent, requester.DefaultUserAgent, "浏览器标识"},
		[]string{"pcs_ua", c.PCSUA, "", "PCS 浏览器标识"},
		[]string{"pcs_addr", c.PCSAddr, "pcs.baidu.com", "PCS 服务器地址"},
		[]string{"fix_pcs_addr", fmt.Sprint(c.FixPCSAddr), "false", "不使用动态PCS服务器地址, 通常情况保持默认即可"},
		[]string{"pan_ua", c.PanUA, baidupcs.NetdiskUA, "Pan 浏览器标识"},
		[]string{"proxy", c.Proxy, "", "设置代理, 支持 http/socks5 代理"},
		[]string{"proxy_hostnames", c.ProxyHostnames, "", "设置走代理的域名范围, 多个域名以逗号分隔, 留空表示全部代理. 国外VPS遇上传问题可尝试代理pan.baidu.com回国"},
		[]string{"local_addrs", c.LocalAddrs, "", "设置本地网卡地址, 多个地址用逗号隔开"},
	})
	tb.Render()
}
