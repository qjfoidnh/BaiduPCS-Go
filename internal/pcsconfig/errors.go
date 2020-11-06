package pcsconfig

import (
	"errors"
)

var (
	//ErrNotLogin 未登录帐号错误
	ErrNotLogin = errors.New("baidu user not login")
	//ErrConfigFilePathNotSet 未设置配置文件
	ErrConfigFilePathNotSet = errors.New("config file not set")
	//ErrConfigFileNotExist 未设置Config, 未初始化
	ErrConfigFileNotExist = errors.New("config file not exist")
	//ErrConfigFileNoPermission Config文件无权限访问
	ErrConfigFileNoPermission = errors.New("config file permission denied")
	//ErrConfigContentsParseError 解析Config数据错误
	ErrConfigContentsParseError = errors.New("config contents parse error")
)
