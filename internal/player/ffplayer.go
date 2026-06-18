package player

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
//	"time"
)

// FFplayPlayer FFplay 播放器（简化版，不支持远程控制）
type FFplayPlayer struct {
	cmd     *exec.Cmd
	mu      sync.Mutex
	running bool
	state   State
}

// FFplayFactory FFplay 工厂
type FFplayFactory struct{}

func (f *FFplayFactory) Create() (Player, error) {
	return NewFFplayPlayer(), nil
}

func (f *FFplayFactory) IsAvailable() bool {
	_, err := exec.LookPath("ffplay")
	return err == nil
}

func (f *FFplayFactory) Name() string {
	return "ffplay"
}

// NewFFplayPlayer 创建 FFplay 播放器
func NewFFplayPlayer() *FFplayPlayer {
	return &FFplayPlayer{
		state: StateStopped,
	}
}

// FFplay 不支持远程控制，所以大部分方法返回错误
func (p *FFplayPlayer) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.running {
		return nil
	}
	
	p.running = true
	return nil
}

func (p *FFplayPlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
	
	p.running = false
	p.state = StateStopped
	return nil
}

func (p *FFplayPlayer) IsRunning() bool {
	return p.running
}

func (p *FFplayPlayer) Name() string {
	return "FFplay"
}

func (p *FFplayPlayer) Play(url string) error {
	// FFplay 每次播放需要启动新进程
	p.Stop()
	
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 FFplay 失败: %v", err)
	}
	
	p.cmd = cmd
	p.state = StatePlaying
	
	// 等待播放完成
	go func() {
		cmd.Wait()
		p.state = StateStopped
	}()
	
	return nil
}

// FFplay 不支持这些功能
func (p *FFplayPlayer) Pause() error {
	return fmt.Errorf("FFplay 不支持远程暂停")
}

func (p *FFplayPlayer) Resume() error {
	return fmt.Errorf("FFplay 不支持远程控制")
}

func (p *FFplayPlayer) Next() error {
	return fmt.Errorf("FFplay 不支持播放列表")
}

func (p *FFplayPlayer) Previous() error {
	return fmt.Errorf("FFplay 不支持播放列表")
}

func (p *FFplayPlayer) Seek(position int) error {
	return fmt.Errorf("FFplay 不支持跳转")
}

func (p *FFplayPlayer) AddToPlaylist(url string) error {
	return fmt.Errorf("FFplay 不支持播放列表")
}

func (p *FFplayPlayer) ClearPlaylist() error {
	return fmt.Errorf("FFplay 不支持播放列表")
}

func (p *FFplayPlayer) SetVolume(level int) error {
	return fmt.Errorf("FFplay 不支持音量控制")
}

func (p *FFplayPlayer) SetSpeed(speed float64) error {
	return fmt.Errorf("FFplay 不支持速度控制")
}

func (p *FFplayPlayer) SetSubtitle(subtitleFile string) error {
	return fmt.Errorf("FFplay 不支持字幕加载")
}

func (p *FFplayPlayer) SendCommand(cmd string) error {
	return fmt.Errorf("FFplay 不支持远程命令")
}

func (p *FFplayPlayer) GetState() State {
	return p.state
}


// 补全 Player 接口
func (p *FFplayPlayer) GetDuration() (float64, error) {
return 0, fmt.Errorf("FFplay 不支持获取时长")
}

func (p *FFplayPlayer) RemoveFromPlaylist(index int) error {
return fmt.Errorf("FFplay 不支持播放列表")
}

func (p *FFplayPlayer) PlayAtIndex(index int) error {
return fmt.Errorf("FFplay 不支持播放列表")
}

func (p *FFplayPlayer) SetShuffle(enabled bool) error {
return fmt.Errorf("FFplay 不支持随机播放")
}

func (p *FFplayPlayer) SetRepeat(mode RepeatMode) error {
return fmt.Errorf("FFplay 不支持循环模式")
}

func (p *FFplayPlayer) GetVolume() (int, error) {
return 0, fmt.Errorf("FFplay 不支持音量查询")
}

func (p *FFplayPlayer) Mute() error {
return fmt.Errorf("FFplay 不支持静音")
}

func (p *FFplayPlayer) Unmute() error {
return fmt.Errorf("FFplay 不支持取消静音")
}

func (p *FFplayPlayer) SetAudioTrack(trackID int) error {
return fmt.Errorf("FFplay 不支持音轨切换")
}

func (p *FFplayPlayer) SetSubtitleDelay(delay float64) error {
return fmt.Errorf("FFplay 不支持字幕延迟")
}

func (p *FFplayPlayer) GetPosition() (float64, error) {
return 0, fmt.Errorf("FFplay 不支持获取位置")
}

func (p *FFplayPlayer) SetLoop(start, end float64) error {
return fmt.Errorf("FFplay 不支持 AB 循环")
}

func (p *FFplayPlayer) TakeSnapshot() (string, error) {
return "", fmt.Errorf("FFplay 不支持截图")
}

func (p *FFplayPlayer) GetPlaylist() ([]PlaylistItem, error) {
return nil, fmt.Errorf("FFplay 不支持播放列表")
}
