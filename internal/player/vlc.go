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

// VLCPlayer VLC 播放器
type VLCPlayer struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mu      sync.Mutex
	running bool
	state   State
	verbose bool
	
	// 播放状态
	volume     int
	speed      float64
	repeatMode RepeatMode
	shuffle    bool
}

// VLCFactory VLC 工厂
type VLCFactory struct {
	Verbose bool
}

func (f *VLCFactory) Create() (Player, error) {
	return NewVLCPlayer(f.Verbose), nil
}

func (f *VLCFactory) IsAvailable() bool {
	_, err := exec.LookPath("vlc")
	return err == nil
}

func (f *VLCFactory) Name() string {
	return "vlc"
}

// NewVLCPlayer 创建 VLC 播放器
func NewVLCPlayer(verbose ...bool) *VLCPlayer {
	v := false
	if len(verbose) > 0 {
		v = verbose[0]
	}

	return &VLCPlayer{
		state:      StateStopped,
		verbose:    v,
		volume:     100,
		speed:      1.0,
		repeatMode: RepeatOff,
		shuffle:    false,
	}
}

func (p *VLCPlayer) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	args := []string{
		"--intf", "rc",
		"--rc-fake-tty",
		"--quiet",
		"--no-qt-privacy-ask",
		"--no-qt-system-tray",
	}

	cmd := exec.Command("vlc", args...)

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
		return fmt.Errorf("启动 VLC 失败: %v", err)
	}

	p.cmd = cmd
	p.stdin = stdin
	p.running = true

	time.Sleep(2 * time.Second)

	// 设置初始音量和速度
	p.setVolume(p.volume)
	p.setSpeed(p.speed)

	return nil
}

func (p *VLCPlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return nil
	}

	p.sendCmd("shutdown")

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

func (p *VLCPlayer) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

func (p *VLCPlayer) Name() string {
	return "VLC"
}

// Play 播放文件/URL
func (p *VLCPlayer) Play(url string) error {
	if url=="" {
		return p.sendCmd("play")
	}
	p.state = StatePlaying
	return p.sendCmd("add " + url)
}

func (p *VLCPlayer) AddToPlaylist(url string) error {
	if url == "" {
		return fmt.Errorf("url cannot be empty")
	}
	return p.sendCmd("enqueue " + url)
}

// Pause 暂停
func (p *VLCPlayer) Pause() error {
	state := p.state
	if state == StatePlaying {
		p.state = StatePaused
	} else if state == StatePaused {
		p.state = StatePlaying
	}
	return p.sendCmd("pause")
}

// Resume 继续播放
func (p *VLCPlayer) Resume() error {
	p.state = StatePlaying
	return p.sendCmd("play")
}

// Next 下一首
func (p *VLCPlayer) Next() error {
	return p.sendCmd("next")
}

// Previous 上一首
func (p *VLCPlayer) Previous() error {
	return p.sendCmd("prev")
}

// Seek 跳转到指定位置（秒）
func (p *VLCPlayer) Seek(position int) error {
	return p.sendCmd(fmt.Sprintf("seek %d", position))
}

// AddToPlaylist 添加到播放列表
func (p *VLCPlayer) ToPlaylist(url string) error {
	return p.sendCmd("enqueue " + url)
}

// ClearPlaylist 清空播放列表
func (p *VLCPlayer) ClearPlaylist() error {
	return p.sendCmd("clear")
}

// GetPlaylist 获取播放列表
func (p *VLCPlayer) GetPlaylist() ([]PlaylistItem, error) {
	// VLC 的 playlist 命令输出需要解析
	// 这里返回空列表，实际使用中可以通过解析输出获取
	return []PlaylistItem{}, nil
}

// RemoveFromPlaylist 从播放列表移除
func (p *VLCPlayer) RemoveFromPlaylist(index int) error {
	return p.sendCmd(fmt.Sprintf("delete %d", index))
}

// PlayAtIndex 播放指定索引的项
func (p *VLCPlayer) PlayAtIndex(index int) error {
	return p.sendCmd(fmt.Sprintf("goto %d", index))
}

// SetShuffle 设置随机播放
func (p *VLCPlayer) SetShuffle(enabled bool) error {
	p.shuffle = enabled
	if enabled {
		return p.sendCmd("random on")
	}
	return p.sendCmd("random off")
}

// SetRepeat 设置重复模式
func (p *VLCPlayer) SetRepeat(mode RepeatMode) error {
	p.repeatMode = mode
	switch mode {
	case RepeatOff:
		return p.sendCmd("repeat off")
	case RepeatOne:
		return p.sendCmd("repeat one")
	case RepeatAll:
		return p.sendCmd("repeat all")
	}
	return nil
}

// SetVolume 设置音量 (0-255)
func (p *VLCPlayer) SetVolume(level int) error {
	if level < 0 || level > 255 {
		return fmt.Errorf("音量范围 0-255")
	}
	p.volume = level
	p.setVolume(level); return nil
}

// GetVolume 获取音量
func (p *VLCPlayer) GetVolume() (int, error) {
	return p.volume, nil
}

// Mute 静音
func (p *VLCPlayer) Mute() error {
	return p.sendCmd("volume 0")
}

// Unmute 取消静音
func (p *VLCPlayer) Unmute() error {
	p.setVolume(p.volume); return nil
}

// SetAudioTrack 设置音轨
func (p *VLCPlayer) SetAudioTrack(trackID int) error {
	return p.sendCmd(fmt.Sprintf("atrack %d", trackID))
}

// SetSubtitle 设置字幕文件
func (p *VLCPlayer) SetSubtitle(subtitleFile string) error {
	return p.sendCmd(fmt.Sprintf("subfile %s", subtitleFile))
}

// SetSubtitleDelay 设置字幕延迟（秒）
func (p *VLCPlayer) SetSubtitleDelay(delay float64) error {
	return p.sendCmd(fmt.Sprintf("subdelay %f", delay))
}

// GetPosition 获取当前播放位置（秒）
func (p *VLCPlayer) GetPosition() (float64, error) {
	// 可以通过 get_time 命令获取，这里返回 0
	return 0, nil
}

// GetDuration 获取总时长（秒）
func (p *VLCPlayer) GetDuration() (float64, error) {
	// 可以通过 get_length 命令获取，这里返回 0
	return 0, nil
}

// GetState 获取播放状态
func (p *VLCPlayer) GetState() State {
	return p.state
}

// SetSpeed 设置播放速度 (0.5-2.0)
func (p *VLCPlayer) SetSpeed(speed float64) error {
	if speed < 0.5 || speed > 2.0 {
		return fmt.Errorf("速度范围 0.5-2.0")
	}
	p.speed = speed
	p.setSpeed(speed); return nil
}

// SetLoop 设置 AB 循环
func (p *VLCPlayer) SetLoop(start, end float64) error {
	// VLC 不直接支持 AB 循环，通过循环 seek 实现
	go func() {
		for p.running {
			time.Sleep(time.Duration(end-start) * time.Second)
			if p.running {
				p.sendCmd(fmt.Sprintf("seek %f", start))
			}
		}
	}()
	return nil
}

// TakeSnapshot 截图
func (p *VLCPlayer) TakeSnapshot() (string, error) {
	filename := fmt.Sprintf("snapshot_%d.png", time.Now().Unix())
	if err := p.sendCmd("snapshot"); err != nil {
		return "", err
	}
	return filename, nil
}

// SendCommand 发送原始命令
func (p *VLCPlayer) SendCommand(cmd string) error {
	// 解析一些常用命令
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	// 处理复合命令
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
			p.shuffle = parts[1] == "on"
		} else {
			p.shuffle = !p.shuffle
		}
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
	case "shutdown", "quit":
		return p.Stop()
	}

	return p.sendCmd(cmd)
}

// 内部方法

func (p *VLCPlayer) sendCmd(cmd string) error {
	if !p.running || p.stdin == nil {
		return fmt.Errorf("播放器未运行")
	}

	cmd = strings.TrimSpace(cmd) + "\n"
	_, err := p.stdin.Write([]byte(cmd))
	return err
}

func (p *VLCPlayer) setVolume(level int) {
	p.sendCmd(fmt.Sprintf("volume %d", level))
}

func (p *VLCPlayer) setSpeed(speed float64) {
	p.sendCmd(fmt.Sprintf("rate %.2f", speed))
}

