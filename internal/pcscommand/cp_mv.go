package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/baidupcs/pcserror"
	"path"
	"strings"
)

// RunCopy 执行 批量拷贝文件/目录
func RunCopy(paths ...string) {
	runCpMvOp("copy", paths...)
}

// RunMove 执行 批量 重命名/移动 文件/目录
func RunMove(paths ...string) {
	runCpMvOp("move", paths...)
}

func runCpMvOp(op string, paths ...string) {
	err := cpmvPathValid(paths...) // 检查路径的有效性, 目前只是判断数量
	if err != nil {
		fmt.Printf("%s path error, %s\n", op, err)
		return
	}

	froms, to := cpmvParsePath(paths...) // 分割

	froms, err = matchPathByShellPattern(froms...)
	if err != nil {
		fmt.Println(err)
		return
	}
	to = GetActiveUser().PathJoin(to)

	// 尝试匹配
	if strings.ContainsAny(to, baidupcs.ShellPatternCharacters) {
		tos, _ := matchPathByShellPattern(to)

		switch len(tos) {
		case 0:
			// do nothing
		case 1:
			to = tos[0]
		default:
			fmt.Printf("目标目录有 %d 条匹配结果, 请检查通配符\n", len(tos))
			return
		}
	}

	pcs := GetBaiduPCS()
	toInfo, pcsError := pcs.FilesDirectoriesMeta(to)
	switch {
	case toInfo != nil && toInfo.Path != to:
		fallthrough
	case pcsError != nil && pcsError.GetErrType() == pcserror.ErrTypeRemoteError:
		// 判断路径是否存在
		// 如果不存在, 则为重命名或同目录拷贝操作

		// 如果 froms 数不是1, 则意义不明确.
		if len(froms) != 1 {
			fmt.Println(err)
			return
		}

		if op == "copy" { // 拷贝
			err = pcs.Copy(&baidupcs.CpMvJSON{
				From: froms[0],
				To:   to,
			})
			if err != nil {
				fmt.Println(err)
				fmt.Println("文件/目录拷贝失败: ")
				fmt.Printf("%s <-> %s\n", froms[0], to)
				return
			}
			fmt.Println("文件/目录拷贝成功: ")
			fmt.Printf("%s <-> %s\n", froms[0], to)
		} else { // 重命名
			err = pcs.Rename(froms[0], to)
			if err != nil {
				fmt.Println(err)
				fmt.Println("重命名失败: ")
				fmt.Printf("%s -> %s\n", froms[0], to)
				return
			}
			fmt.Println("重命名成功: ")
			fmt.Printf("%s -> %s\n", froms[0], to)
		}
		return
	case pcsError != nil && pcsError.GetErrType() != pcserror.ErrTypeRemoteError:
		fmt.Println(pcsError)
		return
	}

	if !toInfo.Isdir {
		fmt.Printf("目标 %s 不是一个目录, 操作失败\n", toInfo.Path)
		return
	}

	cj := new(baidupcs.CpMvListJSON)
	cj.List = make([]*baidupcs.CpMvJSON, len(froms))
	for k := range froms {
		cj.List[k] = &baidupcs.CpMvJSON{
			From: froms[k],
			To:   path.Clean(to + baidupcs.PathSeparator + path.Base(froms[k])),
		}
	}

	switch op {
	case "copy":
		err = pcs.Copy(cj.List...)
		if err != nil {
			fmt.Println(err)
			fmt.Println("操作失败, 以下文件/目录拷贝失败: ")
			fmt.Println(cj)
			return
		}
		fmt.Println("操作成功, 以下文件/目录拷贝成功: ")
		fmt.Println(cj)
	case "move":
		err = pcs.Move(cj.List...)
		if err != nil {
			fmt.Println(err)
			fmt.Println("操作失败, 以下文件/目录移动失败: ")
			fmt.Println(cj)
			return
		}
		fmt.Println("操作成功, 以下文件/目录移动成功: ")
		fmt.Println(cj)
	default:
		panic("Unknown operation:" + op)
	}
	return
}

// cpmvPathValid 检查路径的有效性
func cpmvPathValid(paths ...string) (err error) {
	if len(paths) <= 1 {
		return fmt.Errorf("参数不完整")
	}

	return nil
}

// cpmvParsePath 解析路径
func cpmvParsePath(paths ...string) (froms []string, to string) {
	if len(paths) == 0 {
		return nil, ""
	}
	froms = paths[:len(paths)-1]
	to = paths[len(paths)-1]
	return
}
