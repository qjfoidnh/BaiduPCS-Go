package pcscommand

import (
	"errors"
	"fmt"
)

var (
	// ErrShellPatternMultiRes 多条通配符匹配结果
	ErrShellPatternMultiRes = errors.New("多条通配符匹配结果")
	// ErrShellPatternNoHit 未匹配到路径
	ErrShellPatternNoHit = errors.New("未匹配到路径, 请检测通配符")
)

// ListTask 队列状态 (基类)
type ListTask struct {
	ID       int // 任务id
	MaxRetry int // 最大重试次数
	retry    int // 任务失败的重试次数
}

// RunTestShellPattern 执行测试通配符
func RunTestShellPattern(pattern string) {
	pcs := GetBaiduPCS()
	paths, err := pcs.MatchPathByShellPattern(GetActiveUser().PathJoin(pattern))
	if err != nil {
		fmt.Println(err)
		return
	}
	for k := range paths {
		fmt.Printf("%s\n", paths[k])
	}
	return
}

func matchPathByShellPatternOnce(pattern *string) error {
	paths, err := GetBaiduPCS().MatchPathByShellPattern(GetActiveUser().PathJoin(*pattern))
	if err != nil {
		return err
	}
	switch len(paths) {
	case 0:
		return ErrShellPatternNoHit
	case 1:
		*pattern = paths[0]
	default:
		return ErrShellPatternMultiRes
	}

	return nil
}

func matchPathByShellPattern(patterns ...string) (pcspaths []string, err error) {
	acUser, pcs := GetActiveUser(), GetBaiduPCS()
	for k := range patterns {
		ps, err := pcs.MatchPathByShellPattern(acUser.PathJoin(patterns[k]))
		if err != nil {
			return nil, err
		}

		pcspaths = append(pcspaths, ps...)
	}
	return pcspaths, nil
}
