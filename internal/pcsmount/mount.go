//go:build windows || cgo

package pcsmount

import (
	"fmt"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/winfsp/cgofuse/fuse"
)

// Mount blocks until the file system is unmounted.
func Mount(pcs *baidupcs.BaiduPCS, mountpoint string, options *Options) (err error) {
	if mountpoint == "" {
		return fmt.Errorf("挂载点不能为空")
	}
	if options == nil {
		options = &Options{}
	}
	fs, err := NewFileSystem(pcs, options.RemoteRoot, options.CacheTTL)
	if err != nil {
		return err
	}

	host := fuse.NewFileSystemHost(fs)
	host.SetCapReaddirPlus(true)
	host.SetUseIno(true)

	args := []string{"-o", "ro"}
	if options.Debug {
		args = append(args, "-d")
	}
	if options.SingleThread {
		args = append(args, "-s")
	}
	for _, option := range options.FuseOptions {
		if option != "" {
			args = append(args, "-o", option)
		}
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("启动 FUSE 失败: %v", recovered)
		}
	}()
	if !host.Mount(mountpoint, args) {
		return fmt.Errorf("FUSE 挂载失败")
	}
	return nil
}
