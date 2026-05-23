package player

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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
}

// MPVFactory MPV 工厂
type MPVFactory struct{}

func (f *MPVFactory) Create() (Player, error) {
	return NewMPVPlayer(), nil
}

func (f *MPVFactory) IsAvailable() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

func (f *MPVFactory) Name() string {
	return "mpv"
}

// NewMPVPlayer 创建 MPV 播放器
func NewMPVPlayer() *MPVPlayer {
	return &MPVPlayer{
		state: StateStopped,
	}
}

func (p *MPVPlayer) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	// MPV 使用 --input-ipc-server 或者 --input-file
	// 这里使用 stdin 方式
	cmd := exec.Command("mpv", "--idle=yes", "--input-terminal=yes", "--terminal=yes")
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %v", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 MPV 失败: %v", err)
	}

	p.cmd = cmd
	p.stdin = stdin
	p.running = true

	time.Sleep(2 * time.Second)
	
	fmt.Println("MPV 已启动")
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
	fmt.Println("MPV 已停止")
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

func (p *MPVPlayer) SetVolume(level int) error {
	// MPV 音量范围 0-100
	mpvVolume := level * 100 / 255
	return p.sendCmd(fmt.Sprintf("set volume %d", mpvVolume))
}

func (p *MPVPlayer) SetSpeed(speed float64) error {
	return p.sendCmd(fmt.Sprintf("set speed %.2f", speed))
}

func (p *MPVPlayer) SetSubtitle(subtitleFile string) error {
	return p.sendCmd(fmt.Sprintf("sub-add \"%s\"", subtitleFile))
}

func (p *MPVPlayer) SendCommand(cmd string) error {
	return p.sendCmd(cmd)
}

func (p *MPVPlayer) GetState() State {
	return p.state
}

func (p *MPVPlayer) sendCmd(cmd string) error {
	if !p.running || p.stdin == nil {
		return fmt.Errorf("播放器未运行")
	}
	
	cmd = strings.TrimSpace(cmd) + "\n"
	_, err := p.stdin.Write([]byte(cmd))
	return err
}

