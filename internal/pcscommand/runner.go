package pcscommand

import (
	"io"
	"os"
)

var (
	// DefaultRunner 默认 Runner
	DefaultRunner = Runner{
		Output: os.Stdout,
	}
)

type (
	// Runner 执行器
	Runner struct {
		Output       io.Writer
		IsBackground bool
	}
)
