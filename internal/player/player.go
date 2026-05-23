package player

import (
//	"fmt"
//	"io"
)

// State 播放状态
type State int

const (
	StateStopped State = iota
	StatePlaying
	StatePaused
	StateBuffering
	StateError
)

func (s State) String() string {
	switch s {
	case StateStopped:
		return "已停止"
	case StatePlaying:
		return "播放中"
	case StatePaused:
		return "已暂停"
	case StateBuffering:
		return "缓冲中"
	case StateError:
		return "错误"
	default:
		return "未知"
	}
}

// Player 播放器接口
type Player interface {
	// 生命周期
	Start() error
	Stop() error
	IsRunning() bool
	Name() string
	
	// 播放控制
	Play(url string) error
	Pause() error
	Resume() error
	Next() error
	Previous() error
	Seek(position int) error
	
	// 播放列表
	AddToPlaylist(url string) error
	ClearPlaylist() error
	
	// 音量控制
	SetVolume(level int) error
	
	// 其他
	SetSpeed(speed float64) error
	SetSubtitle(subtitleFile string) error
	SendCommand(cmd string) error
	
	// 获取状态
	GetState() State
}

// PlayerFactory 播放器工厂接口
type PlayerFactory interface {
	Create() (Player, error)
	IsAvailable() bool
	Name() string
}

