// Package pcsconfig 配置包
package pcsconfig

import (
	"github.com/json-iterator/go"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/jsonhelper"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsverbose"
	"github.com/qjfoidnh/BaiduPCS-Go/requester"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	// EnvConfigDir 配置路径环境变量
	EnvConfigDir = "BAIDUPCS_GO_CONFIG_DIR"
	// ConfigName 配置文件名
	ConfigName = "pcs_config.json"
)

var (
	pcsConfigVerbose = pcsverbose.New("PCSCONFIG")
	configFilePath   = filepath.Join(GetConfigDir(), ConfigName)

	// Config 配置信息, 由外部调用
	Config = NewConfig(configFilePath)
)

// PCSConfig 配置详情
type PCSConfig struct {
	BaiduActiveUID uint64        `json:"baidu_active_uid"`
	BaiduUserList  BaiduUserList `json:"baidu_user_list"`

	AppID int `json:"appid"` // appid

	CacheSize         int `json:"cache_size"`          // 下载缓存
	MaxParallel       int `json:"max_parallel"`        // 最大下载并发量
	MaxUploadParallel int `json:"max_upload_parallel"` // 最大上传并发量
	MaxDownloadLoad   int `json:"max_download_load"`   // 同时进行下载文件的最大数量
	MaxUploadLoad     int `json:"max_upload_load"`     // 同时进行上传文件的最大数量

	MaxDownloadRate int64 `json:"max_download_rate"` // 限制最大下载速度
	MaxUploadRate   int64 `json:"max_upload_rate"`   // 限制最大上传速度

	UserAgent      string `json:"user_agent"`           // 浏览器标识
	PCSUA          string `json:"pcs_ua"`               // PCS浏览器标识
	PCSAddr        string `json:"pcs_addr"`             // PCS服务器域名
	PanUA          string `json:"pan_ua"`               // PAN浏览器标识
	SaveDir        string `json:"savedir"`              // 下载储存路径
	EnableHTTPS    bool   `json:"enable_https"`         // 启用https
	FixPCSAddr     bool   `json:"fix_pcs_addr"`         //上传不使用动态PCS服务器域名
	ForceLogin     string `json:"force_login_username"` // 强制登录
	Proxy          string `json:"proxy"`                // 代理
	ProxyHostnames string `json:"proxy_hostnames"`      // 走代理的域名范围
	LocalAddrs     string `json:"local_addrs"`          // 本地网卡地址
	NoCheck        bool   `json:"no_check"`             // 禁用下载md5校验
	IgnoreIllegal  bool   `json:"ignore_illegal"`       // 禁用上传文件名非法字符检查
	UPolicy        string `json:"u_policy"`             // 上传重名文件处理策略

	configFilePath string
	configFile     *os.File
	fileMu         sync.Mutex
	activeUser     *Baidu
	pcs            *baidupcs.BaiduPCS
}

// NewConfig 返回 PCSConfig 指针对象
func NewConfig(configFilePath string) *PCSConfig {
	c := &PCSConfig{
		configFilePath: configFilePath,
	}
	return c
}

// Init 初始化配置
func (c *PCSConfig) Init() error {
	return c.init()
}

// Reload 从文件重载配置
func (c *PCSConfig) Reload() error {
	return c.init()
}

// Close 关闭配置文件
func (c *PCSConfig) Close() error {
	if c.configFile != nil {
		err := c.configFile.Close()
		c.configFile = nil
		return err
	}
	return nil
}

// Save 保存配置信息到配置文件
func (c *PCSConfig) Save() error {
	// 检测配置项是否合法, 不合法则自动修复
	c.fix()

	err := c.lazyOpenConfigFile()
	if err != nil {
		return err
	}

	c.fileMu.Lock()
	defer c.fileMu.Unlock()

	data, err := jsoniter.MarshalIndent(c, "", " ")
	if err != nil {
		// json数据生成失败
		panic(err)
	}

	// 减掉多余的部分
	err = c.configFile.Truncate(int64(len(data)))
	if err != nil {
		return err
	}

	_, err = c.configFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	_, err = c.configFile.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *PCSConfig) init() error {
	if c.configFilePath == "" {
		return ErrConfigFileNotExist
	}

	c.InitDefaultConfig()
	err := c.loadConfigFromFile()
	if err != nil {
		return err
	}

	// 载入配置
	// 如果 activeUser 已初始化, 则跳过
	if c.activeUser != nil && c.activeUser.UID == c.BaiduActiveUID {
		return nil
	}

	c.activeUser, err = c.GetBaiduUser(&BaiduBase{
		UID: c.BaiduActiveUID,
	})
	if err != nil {
		return err
	}
	c.pcs = c.activeUser.BaiduPCS()
	c.pcs.SetPCSAddr(c.PCSAddr)

	// 设置全局User-Agent
	requester.UserAgent = c.UserAgent
	// 设置全局代理
	requester.SetGlobalProxy(c.Proxy)
	// 设置代理规则
	requester.SetProxyHostnameRules(c.ProxyHostnames)
	// 设置本地网卡地址
	requester.SetLocalTCPAddrList(strings.Split(c.LocalAddrs, ",")...)

	return nil
}

// lazyOpenConfigFile 打开配置文件
func (c *PCSConfig) lazyOpenConfigFile() (err error) {
	if c.configFile != nil {
		return nil
	}

	c.fileMu.Lock()
	os.MkdirAll(filepath.Dir(c.configFilePath), 0700)
	c.configFile, err = os.OpenFile(c.configFilePath, os.O_CREATE|os.O_RDWR, 0600)
	c.fileMu.Unlock()

	if err != nil {
		if os.IsPermission(err) {
			return ErrConfigFileNoPermission
		}
		if os.IsExist(err) {
			return ErrConfigFileNotExist
		}
		return err
	}
	return nil
}

// loadConfigFromFile 载入配置
func (c *PCSConfig) loadConfigFromFile() (err error) {
	err = c.lazyOpenConfigFile()
	if err != nil {
		return err
	}

	// 未初始化
	info, err := c.configFile.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		err = c.Save()
		return err
	}

	c.fileMu.Lock()
	defer c.fileMu.Unlock()

	_, err = c.configFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	err = jsonhelper.UnmarshalData(c.configFile, c)
	if err != nil {
		return ErrConfigContentsParseError
	}
	return nil
}

func (c *PCSConfig) InitDefaultConfig() {
	c.AppID = 266719
	c.CacheSize = 65536
	c.MaxParallel = 1
	c.MaxUploadParallel = 4
	c.MaxUploadLoad = 4
	c.MaxDownloadLoad = 1
	c.UserAgent = requester.UserAgent
	c.PCSUA = ""
	c.PCSAddr = "pcs.baidu.com"
	c.PanUA = baidupcs.NetdiskUA
	c.EnableHTTPS = true
	c.NoCheck = true
	c.UPolicy = baidupcs.SkipPolicy
	c.Proxy = ""
	c.LocalAddrs = ""
	c.IgnoreIllegal = true
	c.ForceLogin = ""
	c.EnableHTTPS = true

	// 设置默认的下载路径
	switch runtime.GOOS {
	case "windows":
		c.SaveDir = pcsutil.ExecutablePathJoin("Downloads")
	case "android":
		// TODO: 获取完整的的下载路径
		c.SaveDir = "/sdcard/Download"
	default:
		dataPath, ok := os.LookupEnv("HOME")
		if !ok {
			pcsConfigVerbose.Warn("Environment HOME not set")
			c.SaveDir = pcsutil.ExecutablePathJoin("Downloads")
		} else {
			c.SaveDir = filepath.Join(dataPath, "Downloads")
		}
	}
}

// GetConfigDir 获取配置路径
func GetConfigDir() string {
	// 从环境变量读取
	configDir, ok := os.LookupEnv(EnvConfigDir)
	if ok {
		if filepath.IsAbs(configDir) {
			return configDir
		}
		// 如果不是绝对路径, 从程序目录寻找
		return pcsutil.ExecutablePathJoin(configDir)
	}

	// 使用旧版
	// 如果旧版的配置文件存在, 则使用旧版
	oldConfigDir := pcsutil.ExecutablePath()
	_, err := os.Stat(filepath.Join(oldConfigDir, ConfigName))
	if err == nil {
		return oldConfigDir
	}

	switch runtime.GOOS {
	case "windows":
		dataPath, ok := os.LookupEnv("APPDATA")
		if !ok {
			pcsConfigVerbose.Warn("Environment APPDATA not set")
			return oldConfigDir
		}
		return filepath.Join(dataPath, "BaiduPCS-Go")
	default:
		dataPath, ok := os.LookupEnv("HOME")
		if !ok {
			pcsConfigVerbose.Warn("Environment HOME not set")
			return oldConfigDir
		}
		configDir = filepath.Join(dataPath, ".config", "BaiduPCS-Go")

		// 检测是否可写
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			pcsConfigVerbose.Warnf("check config dir error: %s\n", err)
			return oldConfigDir
		}
		return configDir
	}
}

func (c *PCSConfig) fix() {
	if c.CacheSize < 1024 {
		c.CacheSize = 1024
	}
	if c.MaxParallel < 1 {
		c.MaxParallel = 1
	}
	if c.MaxUploadParallel < 1 {
		c.MaxUploadParallel = 1
	}
	if c.MaxDownloadLoad < 1 {
		c.MaxDownloadLoad = 1
	}
	if c.MaxUploadLoad < 1 {
		c.MaxUploadLoad = 1
	}
	if c.UPolicy != baidupcs.SkipPolicy && c.UPolicy != baidupcs.OverWritePolicy && c.UPolicy != baidupcs.RsyncPolicy {
		c.UPolicy = baidupcs.SkipPolicy
	}
}
