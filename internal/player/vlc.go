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

// VLCPlayer VLC 播放器
type VLCPlayer struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mu      sync.Mutex
	running bool
	state   State
}

// VLCFactory VLC 工厂
type VLCFactory struct{}

func (f *VLCFactory) Create() (Player, error) {
	return NewVLCPlayer(), nil
}

func (f *VLCFactory) IsAvailable() bool {
	_, err := exec.LookPath("vlc")
	return err == nil
}

func (f *VLCFactory) Name() string {
	return "vlc"
}

// NewVLCPlayer 创建 VLC 播放器
func NewVLCPlayer() *VLCPlayer {
	return &VLCPlayer{
		state: StateStopped,
	}
}

func (p *VLCPlayer) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	cmd := exec.Command("vlc", "--intf", "rc", "--rc-fake-tty","--quiet")
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %v", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 VLC 失败: %v", err)
	}

	p.cmd = cmd
	p.stdin = stdin
	p.running = true

	time.Sleep(2 * time.Second)
	
	fmt.Println("VLC 已启动")
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
	fmt.Println("VLC 已停止")
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

func (p *VLCPlayer) Play(url string) error {
	return p.sendCmd("add " + url)
}

func (p *VLCPlayer) Pause() error {
	p.state = StatePaused
	return p.sendCmd("pause")
}

func (p *VLCPlayer) Resume() error {
	p.state = StatePlaying
	return p.sendCmd("play")
}

func (p *VLCPlayer) Next() error {
	return p.sendCmd("next")
}

func (p *VLCPlayer) Previous() error {
	return p.sendCmd("prev")
}

func (p *VLCPlayer) Seek(position int) error {
	return p.sendCmd(fmt.Sprintf("seek %d", position))
}

func (p *VLCPlayer) AddToPlaylist(url string) error {
	return p.sendCmd("enqueue " + url)
}

func (p *VLCPlayer) ClearPlaylist() error {
	return p.sendCmd("clear")
}

func (p *VLCPlayer) SetVolume(level int) error {
	return p.sendCmd(fmt.Sprintf("volume %d", level))
}

func (p *VLCPlayer) SetSpeed(speed float64) error {
	return p.sendCmd(fmt.Sprintf("rate %.2f", speed))
}

func (p *VLCPlayer) SetSubtitle(subtitleFile string) error {
	return p.sendCmd(fmt.Sprintf("subfile %s", subtitleFile))
}

func (p *VLCPlayer) SendCommand(cmd string) error {
	return p.sendCmd(cmd)
}

func (p *VLCPlayer) GetState() State {
	return p.state
}

func (p *VLCPlayer) sendCmd(cmd string) error {
	if !p.running || p.stdin == nil {
		return fmt.Errorf("播放器未运行")
	}
	
	cmd = strings.TrimSpace(cmd) + "\n"
	_, err := p.stdin.Write([]byte(cmd))
	return err
}

