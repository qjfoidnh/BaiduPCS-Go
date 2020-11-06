package baidulogin

import (
	"github.com/GeertJohan/go.rice"
	"github.com/astaxie/beego/session"
)

var (
	templateFilesBox *rice.Box
	libFilesBox      *rice.Box

	globalSessions *session.Manager // 全局 sessions 管理器
)
