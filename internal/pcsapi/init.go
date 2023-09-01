package pcsapi

import (
	"fmt"

	middleware_auth "github.com/qjfoidnh/BaiduPCS-Go/internal/pcsapi/middleware"
)

// TODO: api的配置和初始化
func Init_api(port int, auth bool, username string, password string) {
	router := middleware_auth.Router
	initRunListDirectory(router)
	initRunChangeDirectory(router)
	initRunSearchFiles(router)
	initRunPWD(router)
	initRunGetMeta(router)
	initRunRemove(router)
	initRunMKdir(router)
	initRunCpMv(router)
	initRunDownload(router)
	initRunUpload(router)
	initRunTransfer(router)
	initRunOffLineDownload(router)
	initRunRecycle(router)
	initRunConfigSet(router)
	middleware_auth.Engine.Run(fmt.Sprintf(":%d", port))
}
