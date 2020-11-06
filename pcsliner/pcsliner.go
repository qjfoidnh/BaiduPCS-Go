// Package pcsliner 命令行 readline 工具包
package pcsliner

import (
	"github.com/peterh/liner"
)

// PCSLiner 封装 *liner.State, 提供更简便的操作
type PCSLiner struct {
	State   *liner.State
	History *LineHistory

	tmode liner.ModeApplier
	lmode liner.ModeApplier

	paused bool
}

// NewLiner 返回 *PCSLiner, 默认设置允许 Ctrl+C 结束
func NewLiner() *PCSLiner {
	pl := &PCSLiner{}
	pl.tmode, _ = liner.TerminalMode()

	line := liner.NewLiner()
	pl.lmode, _ = liner.TerminalMode()

	line.SetMultiLineMode(true)
	line.SetCtrlCAborts(true)

	pl.State = line

	return pl
}

// Pause 暂停服务
func (pl *PCSLiner) Pause() error {
	if pl.paused {
		panic("PCSLiner already paused")
	}

	pl.paused = true
	pl.DoWriteHistory()

	return pl.tmode.ApplyMode()
}

// Resume 恢复服务
func (pl *PCSLiner) Resume() error {
	if !pl.paused {
		panic("PCSLiner is not paused")
	}

	pl.paused = false

	return pl.lmode.ApplyMode()
}

// Close 关闭服务
func (pl *PCSLiner) Close() (err error) {
	err = pl.State.Close()
	if err != nil {
		return err
	}

	if pl.History != nil && pl.History.historyFile != nil {
		return pl.History.historyFile.Close()
	}

	return nil
}
