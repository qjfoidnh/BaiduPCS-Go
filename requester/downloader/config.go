package downloader

import (
	"github.com/iikira/BaiduPCS-Go/requester/transfer"
)

const (
	//CacheSize 默认的下载缓存
	CacheSize = 8192
)

var (
	// MinParallelSize 单个线程最小的数据量
	MinParallelSize int64 = 128 * 1024 // 128kb
)

//Config 下载配置
type Config struct {
	Mode                       transfer.RangeGenMode      // 下载Range分配模式
	MaxParallel                int                        // 最大下载并发量
	CacheSize                  int                        // 下载缓冲
	BlockSize                  int64                      // 每个Range区块的大小, RangeGenMode 为 RangeGenMode2 时才有效
	MaxRate                    int64                      // 限制最大下载速度
	InstanceStateStorageFormat InstanceStateStorageFormat // 断点续传储存类型
	InstanceStatePath          string                     // 断点续传信息路径
	IsTest                     bool                       // 是否测试下载
	TryHTTP                    bool                       // 是否尝试使用 http 连接
}

//NewConfig 返回默认配置
func NewConfig() *Config {
	return &Config{
		MaxParallel: 5,
		CacheSize:   CacheSize,
		IsTest:      false,
	}
}

//Fix 修复配置信息, 使其合法
func (cfg *Config) Fix() {
	fixCacheSize(&cfg.CacheSize)
	if cfg.MaxParallel < 1 {
		cfg.MaxParallel = 1
	}
}

//Copy 拷贝新的配置
func (cfg *Config) Copy() *Config {
	newCfg := *cfg
	return &newCfg
}
