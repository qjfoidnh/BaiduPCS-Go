// Package pcsdebug 调试包
package pcsdebug

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
)

//StartCPUProfile 收集cpu信息
func StartCPUProfile(ctx context.Context, cpuProfile string) {
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not create cpu profile output file: %s", err)
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Can not start cpu profile: %s", err)
			f.Close()
			return
		}
		defer pprof.StopCPUProfile()
	}
	select {
	case <-ctx.Done():
		return
	}
}
