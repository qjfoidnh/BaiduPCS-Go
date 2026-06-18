package player

import (
	"fmt"
	"sort"
)

// ManagerConfig 管理器配置
type ManagerConfig struct {
	AutoDetect bool
	Preferred  []string
	Fallback   string
}

// Manager 播放器管理器
type Manager struct {
	factories map[string]PlayerFactory
	player    Player
	config    ManagerConfig
}

// NewManager 创建播放器管理器
func NewManager(config ManagerConfig) *Manager {
	m := &Manager{
		factories: make(map[string]PlayerFactory),
		config:    config,
	}

	// 注册内置播放器
	m.Register(&VLCFactory{})
	m.Register(&MPVFactory{})

	return m
}

// Register 注册播放器工厂
func (m *Manager) Register(factory PlayerFactory) {
	m.factories[factory.Name()] = factory
}

// DetectPlayers 检测可用播放器
func (m *Manager) DetectPlayers() []string {
	var available []string
	for name, factory := range m.factories {
		if factory.IsAvailable() {
			available = append(available, name)
		}
	}

	sort.Slice(available, func(i, j int) bool {
		pi := m.getPriority(available[i])
		pj := m.getPriority(available[j])
		return pi < pj
	})

	return available
}

func (m *Manager) getPriority(name string) int {
	for i, p := range m.config.Preferred {
		if p == name {
			return i
		}
	}
	return 999
}

// SetPlayer 设置播放器
func (m *Manager) SetPlayer(name string) error {
	if m.player != nil {
		m.player.Stop()
	}

	factory, ok := m.factories[name]
	if !ok {
		return fmt.Errorf("未知播放器: %s", name)
	}

	if !factory.IsAvailable() {
		return fmt.Errorf("播放器不可用: %s", name)
	}

	player, err := factory.Create()
	if err != nil {
		return fmt.Errorf("创建播放器失败: %v", err)
	}
	player.Start()
	m.player = player
	fmt.Printf("已选择播放器: %s\n", name)
	return nil
}

// AutoSelect 自动选择播放器
func (m *Manager) AutoSelect() error {
	available := m.DetectPlayers()

	if len(available) == 0 {
		return fmt.Errorf("没有可用的播放器")
	}

	for _, preferred := range m.config.Preferred {
		for _, avail := range available {
			if avail == preferred {
				return m.SetPlayer(avail)
			}
		}
	}

	return m.SetPlayer(available[0])
}

// GetPlayer 获取当前播放器
func (m *Manager) GetPlayer() Player {
	return m.player
}

// Start 启动播放器
func (m *Manager) Start() error {
	if m.player == nil {
		if m.config.AutoDetect {
			if err := m.AutoSelect(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("未选择播放器")
		}
	}

	return m.player.Start()
}

// Stop 停止播放器
func (m *Manager) Stop() error {
	if m.player != nil {
		return m.player.Stop()
	}
	return nil
}

// Play 播放
func (m *Manager) Play(url string) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.Play(url)
}

// Pause 暂停
func (m *Manager) Pause() error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.Pause()
}

// Resume 继续
func (m *Manager) Resume() error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.Resume()
}

// Next 下一首
func (m *Manager) Next() error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.Next()
}

// Previous 上一首
func (m *Manager) Previous() error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.Previous()
}

// Seek 跳转
func (m *Manager) Seek(position int) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.Seek(position)
}

// AddToPlaylist 添加到播放列表
func (m *Manager) AddToPlaylist(url string) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.AddToPlaylist(url)
}

// ClearPlaylist 清空播放列表
func (m *Manager) ClearPlaylist() error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.ClearPlaylist()
}

// SetVolume 设置音量
func (m *Manager) SetVolume(level int) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SetVolume(level)
}

// SetSpeed 设置速度
func (m *Manager) SetSpeed(speed float64) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SetSpeed(speed)
}

// SetLoop 设置循环
func (m *Manager) SetLoop(start, end float64) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SetLoop(start, end)
}

// SetShuffle 设置随机播放
func (m *Manager) SetShuffle(enabled bool) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SetShuffle(enabled)
}

// SetRepeat 设置重复模式
func (m *Manager) SetRepeat(mode RepeatMode) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SetRepeat(mode)
}

// SendCommand 发送命令
func (m *Manager) SendCommand(cmd string) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SendCommand(cmd)
}

// ListPlayers 列出所有播放器
func (m *Manager) ListPlayers() string {
	var result string
	result += "可用播放器:\n"

	for name, factory := range m.factories {
		status := "✗"
		if factory.IsAvailable() {
			status = "✓"
		}

		current := ""
		if m.player != nil && m.player.Name() == name {
			current = " (当前)"
		}

		result += fmt.Sprintf("  [%s] %s%s\n", status, name, current)
	}

	return result
}

// ensurePlayer 确保播放器可用
func (m *Manager) ensurePlayer() error {
	if m.player == nil {
		if err := m.Start(); err != nil {
			return err
		}
	}

	if !m.player.IsRunning() {
		return m.player.Start()
	}

	return nil
}

