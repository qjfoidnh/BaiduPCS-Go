package player

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MPVPlayer MPV 播放器
type MPVPlayer struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mu      sync.Mutex
	running bool
	state   State
	verbose bool
	
	volume     int
	speed      float64
	repeatMode RepeatMode
	shuffle    bool
}

// MPVFactory MPV 工厂
type MPVFactory struct {
	Verbose bool
}

func (f *MPVFactory) Create() (Player, error) {
	return NewMPVPlayer(f.Verbose), nil
}

func (f *MPVFactory) IsAvailable() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

func (f *MPVFactory) Name() string {
	return "mpv"
}

// NewMPVPlayer 创建 MPV 播放器
func NewMPVPlayer(verbose ...bool) *MPVPlayer {
	v := false
	if len(verbose) > 0 {
		v = verbose[0]
	}

	return &MPVPlayer{
		state:      StateStopped,
		verbose:    v,
		volume:     100,
		speed:      1.0,
		repeatMode: RepeatOff,
		shuffle:    false,
	}
}

func (p *MPVPlayer) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	args := []string{
		"--idle=yes",
		"--input-terminal=yes",
		"--terminal=yes",
		"--no-video",  // 默认只播放音频，根据需要修改
	}

	cmd := exec.Command("mpv", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %v", err)
	}

	if p.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 MPV 失败: %v", err)
	}

	p.cmd = cmd
	p.stdin = stdin
	p.running = true

	time.Sleep(1 * time.Second)

	p.setVolume(p.volume)
	p.setSpeed(p.speed)

	return nil
}

func (p *MPVPlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return nil
	}

	p.sendCmd("quit")

	if p.stdin != nil {
		p.stdin.Close()
		p.stdin = nil
	}

	if p.cmd != nil && p.cmd.Process != nil {
		time.Sleep(500 * time.Millisecond)
		p.cmd.Process.Kill()
		p.cmd = nil
	}

	p.running = false
	p.state = StateStopped

	return nil
}

func (p *MPVPlayer) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

func (p *MPVPlayer) Name() string {
	return "MPV"
}

func (p *MPVPlayer) Play(url string) error {
	p.state = StatePlaying
	return p.sendCmd(fmt.Sprintf("loadfile \"%s\"", url))
}

func (p *MPVPlayer) Pause() error {
	p.state = StatePaused
	return p.sendCmd("cycle pause")
}

func (p *MPVPlayer) Resume() error {
	p.state = StatePlaying
	return p.sendCmd("set pause no")
}

func (p *MPVPlayer) Next() error {
	return p.sendCmd("playlist-next")
}

func (p *MPVPlayer) Previous() error {
	return p.sendCmd("playlist-prev")
}

func (p *MPVPlayer) Seek(position int) error {
	return p.sendCmd(fmt.Sprintf("seek %d absolute", position))
}

func (p *MPVPlayer) AddToPlaylist(url string) error {
	return p.sendCmd(fmt.Sprintf("loadfile \"%s\" append", url))
}

func (p *MPVPlayer) ClearPlaylist() error {
	return p.sendCmd("playlist-clear")
}

func (p *MPVPlayer) GetPlaylist() ([]PlaylistItem, error) {
	return []PlaylistItem{}, nil
}

func (p *MPVPlayer) RemoveFromPlaylist(index int) error {
	return p.sendCmd(fmt.Sprintf("playlist-remove %d", index))
}

func (p *MPVPlayer) PlayAtIndex(index int) error {
	return p.sendCmd(fmt.Sprintf("set playlist-pos %d", index))
}

func (p *MPVPlayer) SetShuffle(enabled bool) error {
	p.shuffle = enabled
	val := "no"
	if enabled {
		val = "yes"
	}
	return p.sendCmd(fmt.Sprintf("set shuffle %s", val))
}

func (p *MPVPlayer) SetRepeat(mode RepeatMode) error {
	p.repeatMode = mode
	switch mode {
	case RepeatOff:
		return p.sendCmd("set loop-file no; set loop-playlist no")
	case RepeatOne:
		return p.sendCmd("set loop-file yes; set loop-playlist no")
	case RepeatAll:
		return p.sendCmd("set loop-file no; set loop-playlist yes")
	}
	return nil
}

func (p *MPVPlayer) SetVolume(level int) error {
	if level < 0 || level > 255 {
		return fmt.Errorf("音量范围 0-255")
	}
	p.volume = level
	mpvVolume := level * 100 / 255
	return p.sendCmd(fmt.Sprintf("set volume %d", mpvVolume))
}

func (p *MPVPlayer) GetVolume() (int, error) {
	return p.volume, nil
}

func (p *MPVPlayer) Mute() error {
	return p.sendCmd("set mute yes")
}

func (p *MPVPlayer) Unmute() error {
	return p.sendCmd("set mute no")
}

func (p *MPVPlayer) SetAudioTrack(trackID int) error {
	return p.sendCmd(fmt.Sprintf("set aid %d", trackID))
}

func (p *MPVPlayer) SetSubtitle(subtitleFile string) error {
	return p.sendCmd(fmt.Sprintf("sub-add \"%s\"", subtitleFile))
}

func (p *MPVPlayer) SetSubtitleDelay(delay float64) error {
	return p.sendCmd(fmt.Sprintf("set sub-delay %f", delay))
}

func (p *MPVPlayer) GetPosition() (float64, error) {
	return 0, nil
}

func (p *MPVPlayer) GetDuration() (float64, error) {
	return 0, nil
}

func (p *MPVPlayer) GetState() State {
	return p.state
}

func (p *MPVPlayer) SetSpeed(speed float64) error {
	if speed < 0.5 || speed > 2.0 {
		return fmt.Errorf("速度范围 0.5-2.0")
	}
	p.speed = speed
	p.setSpeed(speed); return nil
}

func (p *MPVPlayer) SetLoop(start, end float64) error {
	// 使用 mpv 的 ab-loop 功能
	return p.sendCmd(fmt.Sprintf("set ab-loop-a %f; set ab-loop-b %f", start, end))
}

func (p *MPVPlayer) TakeSnapshot() (string, error) {
	filename := fmt.Sprintf("snapshot_%d.png", time.Now().Unix())
	if err := p.sendCmd("screenshot"); err != nil {
		return "", err
	}
	return filename, nil
}

func (p *MPVPlayer) SendCommand(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "play":
		if len(parts) > 1 {
			return p.Play(strings.Join(parts[1:], " "))
		}
		return p.Resume()
	case "pause":
		return p.Pause()
	case "stop":
		p.state = StateStopped
	case "next":
		return p.Next()
	case "prev":
		return p.Previous()
	case "volume":
		if len(parts) > 1 {
			if vol, err := strconv.Atoi(parts[1]); err == nil {
				return p.SetVolume(vol)
			}
		}
	case "speed":
		if len(parts) > 1 {
			if spd, err := strconv.ParseFloat(parts[1], 64); err == nil {
				return p.SetSpeed(spd)
			}
		}
	case "seek":
		if len(parts) > 1 {
			if pos, err := strconv.Atoi(parts[1]); err == nil {
				return p.Seek(pos)
			}
		}
	case "random", "shuffle":
		if len(parts) > 1 {
			return p.SetShuffle(parts[1] == "on")
		}
		p.shuffle = !p.shuffle
		return p.SetShuffle(p.shuffle)
	case "repeat":
		if len(parts) > 1 {
			switch parts[1] {
			case "off":
				return p.SetRepeat(RepeatOff)
			case "one":
				return p.SetRepeat(RepeatOne)
			case "all":
				return p.SetRepeat(RepeatAll)
			}
		}
	case "loop":
		if len(parts) >= 3 {
			start, _ := strconv.ParseFloat(parts[1], 64)
			end, _ := strconv.ParseFloat(parts[2], 64)
			return p.SetLoop(start, end)
		}
	case "subtitle", "sub":
		if len(parts) > 1 {
			return p.SetSubtitle(strings.Join(parts[1:], " "))
		}
	case "snapshot", "snap":
		_, err := p.TakeSnapshot()
		return err
	case "mute":
		return p.Mute()
	case "unmute":
		return p.Unmute()
	}

	return p.sendCmd(cmd)
}

func (p *MPVPlayer) sendCmd(cmd string) error {
	if !p.running || p.stdin == nil {
		return fmt.Errorf("播放器未运行")
	}

	cmd = strings.TrimSpace(cmd) + "\n"
	_, err := p.stdin.Write([]byte(cmd))
	return err
}

func (p *MPVPlayer) setVolume(level int) {
	mpvVolume := level * 100 / 255
	p.sendCmd(fmt.Sprintf("set volume %d", mpvVolume))
}

func (p *MPVPlayer) setSpeed(speed float64) {
	p.sendCmd(fmt.Sprintf("set speed %.2f", speed))
}

