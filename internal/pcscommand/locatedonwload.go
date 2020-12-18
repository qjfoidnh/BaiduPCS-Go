package pcscommand

import (
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"os"
	"strconv"
)

// RunLocateDownload 执行获取直链
func RunLocateDownload(pcspaths []string, opt *LocateDownloadOption) {
	if opt == nil {
		opt = &LocateDownloadOption{}
	}

	absPaths, err := matchPathByShellPattern(pcspaths...)
	if err != nil {
		fmt.Println(err)
		return
	}

	pcs := GetBaiduPCS()

	if opt.FromPan {
		fds, err := pcs.FilesDirectoriesBatchMeta(absPaths...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fidList := make([]int64, 0, len(fds))
		for i := range fds {
			fidList = append(fidList, fds[i].FsID)
		}

		list, err := pcs.LocatePanAPIDownload(fidList...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		tb := pcstable.NewTable(os.Stdout)
		tb.SetHeader([]string{"#", "fs_id", "路径", "链接"})

		var (
			i          int
			fidStrList = converter.SliceInt64ToString(fidList)
		)
		for k := range fidStrList {
			for i = range list {
				if fidStrList[k] == list[i].FsID {
					tb.Append([]string{strconv.Itoa(k), list[i].FsID, fds[k].Path, list[i].Dlink})
					list = append(list[:i], list[i+1:]...)
					break
				}
			}
		}
		tb.Render()
		fmt.Printf("\n注意: 以上链接不能直接访问, 需要登录百度帐号才可以下载\n")
		return
	}

	for i, pcspath := range absPaths {
		info, err := pcs.LocateDownload(pcspath)
		if err != nil {
			fmt.Printf("[%d] %s, 路径: %s\n", i, err, pcspath)
			continue
		}

		fmt.Printf("[%d] %s: \n", i, pcspath)
		tb := pcstable.NewTable(os.Stdout)
		tb.SetHeader([]string{"#", "链接"})
		for k, u := range info.URLStrings(pcsconfig.Config.EnableHTTPS) {
			tb.Append([]string{strconv.Itoa(k), u.String()})
		}
		tb.Render()
		fmt.Println()
	}
	fmt.Printf("提示: 访问下载链接, 需将下载器的 User-Agent 设置为: %s\n", pcsconfig.Config.PanUA)
}
