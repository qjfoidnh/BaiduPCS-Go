package baidupcs

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"github.com/tidwall/gjson"
)

// GetShareListBody 请求分享文件列表接口（/share/list），errno==0 时返回响应体。
//
// 供 transfer --list 只列清单不转存使用：调用方须已完成 AccessSharePage 与提取码校验
// （BDCLND 已在 cookie jar 中）。失败时 errMsg 与 ExtractShareInfo 同口径，便于上层沿用
// 转存失败的错误分类。
func (pcs *BaiduPCS) GetShareListBody(listURL, referer string) (body []byte, errMsg string) {
	dataReadCloser, panError := pcs.sendReqReturnReadCloser(reqTypePan, OperationShareFileSavetoLocal, http.MethodGet, listURL, nil, map[string]string{
		"User-Agent":   requester.UserAgent,
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Referer":      referer,
	})
	if panError != nil {
		return nil, "提交分享项查询请求时发生错误"
	}
	defer dataReadCloser.Close()
	body, err := ioutil.ReadAll(dataReadCloser)
	if err != nil {
		return nil, "读取分享项查询响应失败"
	}
	if !gjson.ValidBytes(body) {
		return nil, "返回json解析错误"
	}
	errno := gjson.GetBytes(body, `errno`).Int()
	if errno != 0 {
		// 文案与 ExtractShareInfo 逐字一致（含「错误码-9」），上层按同一套关键词分类
		if errno == 8001 {
			return nil, "已触发验证, 请稍后再试"
		}
		return nil, fmt.Sprintf("未知错误, 错误码%d", errno)
	}
	return body, ""
}
