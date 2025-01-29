package pcsapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
)

type (
	ConfigStructure struct {
		Appid              int    `json:"appid,omitempty" form:"appid,omitempty"`
		CacheSize          string `json:"cache_size,omitempty" form:"cache_size,omitempty"`
		MaxParallel        int    `json:"max_parallel,omitempty" form:"max_parallel,omitempty"`
		MaxUploadParallel  int    `json:"max_upload_parallel,omitempty" form:"max_upload_parallel,omitempty"`
		MaxDownloadLoad    int    `json:"max_download_load,omitempty" form:"max_download_load,omitempty"`
		MaxUploadLoad      int    `json:"max_upload_load,omitempty" form:"max_upload_load,omitempty"`
		MaxDownloadRate    string `json:"max_download_rate,omitempty" form:"max_download_rate,omitempty"`
		MaxUploadRate      string `json:"max_upload_rate,omitempty" form:"max_upload_rate,omitempty"`
		SaveDir            string `json:"save_dir,omitempty" form:"save_dir,omitempty"`
		EnableHttps        bool   `json:"enable_https,omitempty" form:"enable_https,omitempty"`
		IgnoreIllegal      bool   `json:"ignore_illegal,omitempty" form:"ignore_illegal,omitempty"`
		ForceLoginUsername string `json:"force_login_username,omitempty" form:"force_login_username,omitempty"`
		NoCheck            bool   `json:"no_check,omitempty" form:"no_check,omitempty"`
		UploadPolicy       string `json:"upload_policy,omitempty" form:"upload_policy,omitempty"`
		Useragent          string `json:"user_agent,omitempty" form:"user_agent,omitempty"`
		PCS_UA             string `json:"pcs_ua,omitempty" form:"pcs_ua,omitempty"`
		PCS_Addr           string `json:"pcs_addr,omitempty" form:"pcs_addr,omitempty"`
		Pan_UA             string `json:"pan_ua,omitempty" form:"pan_ua,omitempty"`
		Proxy              string `json:"proxy,omitempty" form:"proxy,omitempty"`
		Local_addrs        string `json:"local_addrs,omitempty" form:"local_addrs,omitempty"`
	}
)

func runConfigSet(ctx *gin.Context) {
	args := ConfigStructure{
		Appid:             -1,
		MaxParallel:       -1,
		MaxUploadParallel: -1,
		MaxDownloadLoad:   -1,
		MaxUploadLoad:     -1,
	}
	if err := ctx.ShouldBind(&args); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	if args.Appid != -1 {
		pcsconfig.Config.SetAppID(args.Appid)
	}
	if args.EnableHttps {
		pcsconfig.Config.SetEnableHTTPS(args.EnableHttps)
	}
	if args.IgnoreIllegal {
		pcsconfig.Config.SetIgnoreIllegal(args.IgnoreIllegal)
	}
	if args.ForceLoginUsername != "" {
		pcsconfig.Config.SetForceLogin(args.ForceLoginUsername)
	}
	if args.NoCheck {
		pcsconfig.Config.SetNoCheck(args.NoCheck)
	}
	if args.UploadPolicy != "" {
		pcsconfig.Config.SetUploadPolicy(args.UploadPolicy)
	}
	if args.Useragent != "" {
		pcsconfig.Config.SetUserAgent(args.Useragent)
	}
	if args.PCS_UA != "" {
		pcsconfig.Config.SetPCSUA(args.PCS_UA)
	}
	if args.PCS_Addr != "" {
		match := pcsconfig.Config.SETPCSAddr(args.PCS_Addr)
		if !match {
			err := fmt.Errorf("设置pcs_addr 错误：pcs服务器地址不合法")
			ctx.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
	}
	if args.Pan_UA != "" {
		pcsconfig.Config.SetPanUA(args.Pan_UA)
	}
	if args.CacheSize != "" {
		err := pcsconfig.Config.SetCacheSizeByStr(args.CacheSize)
		if err != nil {
			err = fmt.Errorf("设置cache_size 错误:%s", err)
			ctx.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
	}
	if args.MaxParallel != -1 {
		pcsconfig.Config.MaxParallel = args.MaxParallel
	}
	if args.MaxUploadParallel != -1 {
		pcsconfig.Config.MaxUploadParallel = args.MaxUploadParallel
	}
	if args.MaxDownloadLoad != -1 {
		pcsconfig.Config.MaxDownloadLoad = args.MaxDownloadLoad
	}
	if args.MaxUploadLoad != -1 {
		pcsconfig.Config.MaxUploadLoad = args.MaxUploadLoad
	}
	if args.MaxDownloadRate != "" {
		err := pcsconfig.Config.SetMaxDownloadRateByStr(args.MaxDownloadRate)
		if err != nil {
			err = fmt.Errorf("设置 max_download_rate 错误： %s", err)
			ctx.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
	}
	if args.MaxUploadRate != "" {
		err := pcsconfig.Config.SetMaxUploadRateByStr(args.MaxUploadRate)
		if err != nil {
			err = fmt.Errorf("设置 max_upload_rate 错误： %s", err)
			ctx.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
	}
	if args.SaveDir != "" {
		pcsconfig.Config.SaveDir = args.SaveDir
	}
	if args.Proxy != "" {
		pcsconfig.Config.SetProxy(args.Proxy)
	}
	if args.Local_addrs != "" {
		pcsconfig.Config.SetLocalAddrs(args.Local_addrs)
	}
	err := pcsconfig.Config.Save()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}
	pcsconfig.Config.PrintTable()
	ctx.JSON(http.StatusOK, gin.H{
		"result": "保存配置成功",
	})
	fmt.Printf("\n保存配置成功!\n\n")
}

func initRunConfigSet(router *gin.RouterGroup) {
	g := router.Group("config")
	g.POST("set", runConfigSet)
}
