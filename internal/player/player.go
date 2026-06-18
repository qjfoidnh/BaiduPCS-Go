package player

import (
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
	GetPlaylist() ([]PlaylistItem, error)
	RemoveFromPlaylist(index int) error
	PlayAtIndex(index int) error
	
	// 播放模式
	SetShuffle(enabled bool) error
	SetRepeat(mode RepeatMode) error
	
	// 音量控制
	SetVolume(level int) error
	GetVolume() (int, error)
	Mute() error
	Unmute() error

	// 音轨和字幕
	SetAudioTrack(trackID int) error
	SetSubtitle(subtitleFile string) error
	SetSubtitleDelay(delay float64) error

	// 播放信息
	GetPosition() (float64, error)
	GetDuration() (float64, error)
	GetState() State

	// 高级功能
	SetSpeed(speed float64) error
	SetLoop(start, end float64) error
	TakeSnapshot() (string, error)

	// 原始命令
	SendCommand(cmd string) error
}

// RepeatMode 重复模式
type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatOne
	RepeatAll
)

func (r RepeatMode) String() string {
	switch r {
	case RepeatOff:
		return "关闭"
	case RepeatOne:
		return "单曲循环"
	case RepeatAll:
		return "列表循环"
	default:
		return "未知"
	}
}

// PlaylistItem 播放列表项
type PlaylistItem struct {
	Index    int
	Title    string
	URL      string
	Duration float64
	Type     MediaType
}

// MediaType 媒体类型
type MediaType int

const (
	MediaAudio MediaType = iota
	MediaVideo
	MediaStream
)

func (m MediaType) String() string {
	switch m {
	case MediaAudio:
		return "音频"
	case MediaVideo:
		return "视频"
	case MediaStream:
		return "流媒体"
	default:
		return "未知"
	}
}

// PlayerInfo 播放器信息
type PlayerInfo struct {
	Name        string
	Version     string
	State       State
	Volume      int
	Speed       float64
	RepeatMode  RepeatMode
	Shuffle     bool
	CurrentItem *PlaylistItem
	Position    float64
	Duration    float64
}

// PlayerFactory 播放器工厂接口
type PlayerFactory interface {
	Create() (Player, error)
	IsAvailable() bool
	Name() string
}

