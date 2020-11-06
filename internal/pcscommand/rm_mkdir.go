package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcstable"
	"os"
	"strconv"
)

// RunRemove 执行 批量删除文件/目录
func RunRemove(paths ...string) {
	paths, err := matchPathByShellPattern(paths...)
	if err != nil {
		fmt.Println(err)
		return
	}

	pnt := func() {
		tb := pcstable.NewTable(os.Stdout)
		tb.SetHeader([]string{"#", "文件/目录"})
		for k := range paths {
			tb.Append([]string{strconv.Itoa(k), paths[k]})
		}
		tb.Render()
	}

	err = GetBaiduPCS().Remove(paths...)
	if err != nil {
		fmt.Println(err)
		fmt.Println("操作失败, 以下文件/目录删除失败: ")
		pnt()
		return
	}

	fmt.Println("操作成功, 以下文件/目录已删除, 可在网盘文件回收站找回: ")
	pnt()
}

// RunMkdir 执行 创建目录
func RunMkdir(path string) {
	activeUser := GetActiveUser()
	err := GetBaiduPCS().Mkdir(activeUser.PathJoin(path))
	if err != nil {
		fmt.Printf("创建目录 %s 失败, %s\n", path, err)
		return
	}

	fmt.Println("创建目录成功:", path)
}
