package player

import (
	"fmt"
	"sort"
)

// Manager 播放器管理器
type Manager struct {
	factories map[string]PlayerFactory
	player    Player
	config    ManagerConfig
}

// ManagerConfig 管理器配置
type ManagerConfig struct {
	AutoDetect bool     // 自动检测可用播放器
	Preferred  []string // 首选播放器列表（按优先级排列）
	Fallback   string   // 备选播放器
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
	m.Register(&FFplayFactory{})
	
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
	
	// 按优先级排序
	sort.Slice(available, func(i, j int) bool {
		pi := m.getPriority(available[i])
		pj := m.getPriority(available[j])
		return pi < pj
	})
	
	return available
}

// getPriority 获取播放器优先级
func (m *Manager) getPriority(name string) int {
	for i, p := range m.config.Preferred {
		if p == name {
			return i
		}
	}
	return 999 // 不在首选列表中
}

// SetPlayer 设置播放器
func (m *Manager) SetPlayer(name string) error {
	// 停止当前播放器
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
	
	m.player = player
	fmt.Printf("已选择播放器: %s\n", name)
	return nil
}

// AutoSelect 自动选择最佳播放器
func (m *Manager) AutoSelect() error {
	available := m.DetectPlayers()
	
	if len(available) == 0 {
		return fmt.Errorf("没有可用的播放器")
	}
	
	// 优先选择首选列表中的播放器
	for _, preferred := range m.config.Preferred {
		for _, avail := range available {
			if avail == preferred {
				return m.SetPlayer(avail)
			}
		}
	}
	
	// 没有首选播放器，使用第一个可用的
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

// SendCommand 发送命令
func (m *Manager) SendCommand(cmd string) error {
	if err := m.ensurePlayer(); err != nil {
		return err
	}
	return m.player.SendCommand(cmd)
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

