package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/pcstable"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"github.com/iikira/BaiduPCS-Go/pcsutil/pcstime"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

// RunRecycleList 执行列出回收站文件列表
func RunRecycleList(page int) {
	if page < 1 {
		page = 1
	}

	pcs := GetBaiduPCS()
	fdl, err := pcs.RecycleList(page)
	if err != nil {
		fmt.Println(err)
		return
	}

	tb := pcstable.NewTable(os.Stdout)
	tb.SetHeader([]string{"#", "fs_id", "文件大小", "创建日期", "修改日期", "md5(截图请打码)", "剩余时间", "路径"})
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_LEFT})
	for k, file := range fdl {
		if file.Isdir == 1 {
			tb.Append([]string{strconv.Itoa(k), strconv.FormatInt(file.FsID, 10), "-", pcstime.FormatTime(file.Ctime), pcstime.FormatTime(file.Mtime), file.MD5, strconv.Itoa(file.LeftTime), file.Path + baidupcs.PathSeparator})
			continue
		}
		tb.Append([]string{strconv.Itoa(k), strconv.FormatInt(file.FsID, 10), converter.ConvertFileSize(file.Size, 2), pcstime.FormatTime(file.Ctime), pcstime.FormatTime(file.Mtime), file.MD5, strconv.Itoa(file.LeftTime), file.Path})
	}

	tb.Render()
}

// RunRecycleRestore 执行还原回收站文件或目录
func RunRecycleRestore(fidStrList ...string) {
	var (
		fidList = converter.SliceStringToInt64(fidStrList)
		pcs     = GetBaiduPCS()
		ex, err = pcs.RecycleRestore(fidList...)
	)
	if err != nil {
		fmt.Println(err)
		if len(ex) > 0 {
			fmt.Printf("\n以下的 fs_id 还原成功, 数量: %d\n", len(ex))
			for k := range ex {
				fmt.Println(ex[k].FsID)
			}
		}
		return
	}

	fmt.Printf("还原成功, 数量: %d\n", len(ex))
}

// RunRecycleDelete 执行删除回收站文件或目录
func RunRecycleDelete(fidStrList ...string) {
	var (
		fidList = converter.SliceStringToInt64(fidStrList)
		pcs     = GetBaiduPCS()
		err     = pcs.RecycleDelete(fidList...)
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("删除成功\n")
}

// RunRecycleClear 清空回收站
func RunRecycleClear() {
	pcs := GetBaiduPCS()
	sussNum, err := pcs.RecycleClear()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("清空回收站成功, 数量: %d\n", sussNum)
}
