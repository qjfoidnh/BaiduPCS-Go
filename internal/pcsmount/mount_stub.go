//go:build !windows && !cgo

package pcsmount

import (
	"errors"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

// Mount reports why a Unix build without cgo cannot host FUSE.
func Mount(_ *baidupcs.BaiduPCS, _ string, _ *Options) error {
	return errors.New("当前程序未启用挂载支持；Linux/macOS 需要安装 FUSE 开发库并使用 CGO_ENABLED=1 重新编译")
}
