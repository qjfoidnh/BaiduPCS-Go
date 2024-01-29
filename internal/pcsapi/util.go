package pcsapi

import (
	"errors"

	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
)

var (
	// ErrShellPatternMultiRes 多条通配符匹配结果
	ErrShellPatternMultiRes = errors.New("多条通配符匹配结果")
	// ErrShellPatternNoHit 未匹配到路径
	ErrShellPatternNoHit = errors.New("未匹配到路径, 请检测通配符")
)

// RunTestShellPattern 执行测试通配符
func runTestShellPattern(pattern string) (paths []string, err error) {
	pcs := pcscommand.GetBaiduPCS()
	paths, err = pcs.MatchPathByShellPattern(pcscommand.GetActiveUser().PathJoin(pattern))
	return
}

func matchPathByShellPatternOnce(pattern *string) error {
	paths, err := pcscommand.GetBaiduPCS().MatchPathByShellPattern(pcscommand.GetActiveUser().PathJoin(*pattern))
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
	acUser, pcs := pcscommand.GetActiveUser(), pcscommand.GetBaiduPCS()
	for k := range patterns {
		ps, err := pcs.MatchPathByShellPattern(acUser.PathJoin(patterns[k]))
		if err != nil {
			return nil, err
		}
		pcspaths = append(pcspaths, ps...)
	}
	return pcspaths, nil
}
