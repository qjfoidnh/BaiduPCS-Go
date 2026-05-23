package pcscommand

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/player"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsdownload"
	"github.com/urfave/cli"
)

const (
	cacheExpireHours = 8
)

// URLCache URL 缓存
type URLCache struct {
	mu    sync.RWMutex
	files map[string]CacheItem
}

type CacheItem struct {
	URL      string
	FilePath string
	ExpireAt time.Time
}

var urlCache = &URLCache{
	files: make(map[string]CacheItem),
}

func (c *URLCache) Get(path string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.files[path]; ok {
		if time.Now().Before(item.ExpireAt) {
			return item.URL, true
		}
		delete(c.files, path)
	}
	return "", false
}

func (c *URLCache) Set(path, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.files[path] = CacheItem{
		URL:      url,
		FilePath: path,
		ExpireAt: time.Now().Add(cacheExpireHours * time.Hour),
	}
}

func (c *URLCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files = make(map[string]CacheItem)
}

func (c *URLCache) GetCacheInfo() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("缓存数量: %d\n", len(c.files)))
	
	now := time.Now()
	for path, item := range c.files {
		remaining := item.ExpireAt.Sub(now)
		sb.WriteString(fmt.Sprintf("  %s (剩余 %s)\n", path, remaining.Round(time.Second)))
	}
	
	return sb.String()
}

// 全局播放器管理器
var playerManager *player.Manager

// InitPlayer 初始化播放器
func InitPlayer() {
	config := player.ManagerConfig{
		AutoDetect: true,
		Preferred:  []string{"vlc", "mpv", "ffplay"},
	}
	
	playerManager = player.NewManager(config)
	
	// 自动选择播放器
	if err := playerManager.AutoSelect(); err != nil {
		fmt.Printf("警告: %v\n", err)
		fmt.Println("请安装 VLC 或 MPV 播放器以获得完整体验")
	}
}

// Play 播放入口
func Play(c *cli.Context) error {
	if playerManager == nil {
		InitPlayer()
	}
	
	if c.NArg() == 0 {
		return controller()
	}

	path := c.Args().Get(0)
	return PlayFile(path)
}

// PlayFile 播放文件
func PlayFile(path string) error {
	if playerManager == nil {
		InitPlayer()
	}

	pcs := GetBaiduPCS()
	url, err := getURL(pcs, path)
	if err != nil {
		return fmt.Errorf("获取URL失败: %v", err)
	}

	if err := playerManager.Play(url); err != nil {
		return fmt.Errorf("播放失败: %v", err)
	}
	
	fmt.Printf("开始播放: %s\n", filepath.Base(path))
	return nil
}

// getURL 获取文件下载链接（带缓存）
func getURL(pcs *baidupcs.BaiduPCS, path string) (string, error) {
	if url, ok := urlCache.Get(path); ok {
		fmt.Printf("使用缓存URL: %s\n", path)
		return url, nil
	}

	fmt.Printf("获取新的下载链接: %s\n", path)
	
	dlinks, err := pcsdownload.GetLocateDownloadLinks(pcs, path)
	if err != nil {
		return "", err
	}
	if len(dlinks) == 0 {
		return "", fmt.Errorf("未找到下载链接")
	}

	url := dlinks[0].String()
	urlCache.Set(path, url)
	
	return url, nil
}

// controller 交互式控制器
func controller() error {
	if playerManager == nil {
		InitPlayer()
	}

	p := playerManager.GetPlayer()
	
	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║     BaiduPCS-Go 多媒体播放器             ║")
	fmt.Printf("║     播放器: %-30s║\n", p.Name())
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Println("║  输入 help 查看命令列表                  ║")
	fmt.Println("║  输入 quit 退出控制器                    ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// 确保播放器运行
	if !p.IsRunning() {
		fmt.Printf("正在启动 %s...\n", p.Name())
		if err := playerManager.Start(); err != nil {
			return fmt.Errorf("启动播放器失败: %v", err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nplayer:> ")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Print("player:> ")
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "add", "a":
			handleAdd(parts)
		case "play", "p":
			handlePlay(parts)
		case "pause", "pa":
			handlePause()
		case "stop", "s":
			handleStop()
		case "next", "n":
			handleNext()
		case "prev", "pr":
			handlePrev()
		case "list", "l":
			handleList()
		case "clear", "c":
			handleClear()
		case "status", "st":
			handleStatus()
		case "volume", "v":
			handleVolume(parts)
		case "seek", "sk":
			handleSeek(parts)
		case "speed", "sp":
			handleSpeed(parts)
		case "subtitle", "sub":
			handleSubtitle(parts)
		case "player":
			handlePlayerCommand(parts)
		case "cache":
			handleCache(parts)
		case "help", "h", "?":
			showHelp()
		case"exit","e":
			return nil
		case "quit", "q":
			playerManager.SendCommand("quit")
			return nil
		default:
			// 尝试作为原始命令发送给播放器
			if err := playerManager.SendCommand(line); err != nil {
				fmt.Printf("未知命令: %s (输入 help 查看帮助)\n", cmd)
			}
		}

		fmt.Print("player:> ")
	}

	return scanner.Err()
}

// handlePlayerCommand 播放器管理命令
func handlePlayerCommand(parts []string) {
	if len(parts) < 2 {
		fmt.Println(playerManager.ListPlayers())
		return
	}

	switch parts[1] {
	case "list":
		fmt.Println(playerManager.ListPlayers())
	case "switch":
		if len(parts) < 3 {
			fmt.Println("用法: player switch <播放器名称>")
			fmt.Println(playerManager.ListPlayers())
			return
		}
		if err := playerManager.SetPlayer(parts[2]); err != nil {
			fmt.Printf("切换失败: %v\n", err)
		} else {
			fmt.Printf("已切换到: %s\n", parts[2])
		}
	case "detect":
		available := playerManager.DetectPlayers()
		fmt.Println("检测到的播放器:")
		for _, name := range available {
			fmt.Printf("  ✓ %s\n", name)
		}
	default:
		fmt.Println("用法: player [list|switch|detect]")
	}
}

// handleAdd 处理添加命令
func handleAdd(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: add <百度云路径>")
		return
	}

	path := strings.Join(parts[1:], " ")
	pcs := GetBaiduPCS()
	url, err := getURL(pcs, path)
	if err != nil {
		fmt.Printf("添加失败: %v\n", err)
		return
	}

	if err := playerManager.Play(url); err != nil {
		fmt.Printf("添加失败: %v\n", err)
	} else {
		fmt.Printf("已添加: %s\n", filepath.Base(path))
	}
}

// handlePlay 处理播放命令
func handlePlay(parts []string) {
	if len(parts) < 2 {
		playerManager.SendCommand("play")
		return
	}

	path := strings.Join(parts[1:], " ")
	
	// 直接 URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		if err := playerManager.Play(path); err != nil {
			fmt.Printf("播放失败: %v\n", err)
		} else {
			fmt.Println("开始播放")
		}
		return
	}

	// 百度云路径
	pcs := GetBaiduPCS()
	url, err := getURL(pcs, path)
	if err != nil {
		fmt.Printf("获取链接失败: %v\n", err)
		return
	}

	if err := playerManager.Play(url); err != nil {
		fmt.Printf("播放失败: %v\n", err)
	} else {
		fmt.Printf("开始播放: %s\n", filepath.Base(path))
	}
}

func handlePause() {
	playerManager.SendCommand("pause")
}

func handleStop() {
	playerManager.SendCommand("stop")
}

func handleNext() {
	playerManager.SendCommand("next")
}

func handlePrev() {
	playerManager.SendCommand("prev")
}

func handleList() {
	playerManager.SendCommand("playlist")
}

func handleClear() {
	playerManager.SendCommand("clear")
}

func handleStatus() {
	playerManager.SendCommand("status")
}

func handleVolume(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: volume <0-255>")
		return
	}
	playerManager.SendCommand(fmt.Sprintf("volume %s", parts[1]))
}

func handleSeek(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: seek <秒>")
		return
	}
	playerManager.SendCommand(fmt.Sprintf("seek %s", parts[1]))
}

func handleSpeed(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: speed <0.5-2.0>")
		return
	}
	playerManager.SendCommand(fmt.Sprintf("speed %s", parts[1]))
}

func handleSubtitle(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: subtitle <字幕文件>")
		return
	}
	playerManager.SendCommand(fmt.Sprintf("subtitle %s", parts[1]))
}

func handleCache(parts []string) {
	if len(parts) < 2 {
		fmt.Println("URL 缓存信息:")
		fmt.Print(urlCache.GetCacheInfo())
		return
	}

	switch parts[1] {
	case "clear":
		urlCache.Clear()
		fmt.Println("缓存已清空")
	default:
		fmt.Println("用法: cache [clear]")
	}
}

func showHelp() {
	fmt.Println(`
┌────────────────────────────────────────────────────┐
│                  可用命令列表                       │
├────────────────────────────────────────────────────┤
│ 播放控制:                                          │
│   play/p  [路径]    播放文件/继续播放               │
│   pause/pa          暂停/继续                       │
│   stop/s            停止播放                        │
│   next/n            下一首                          │
│   prev/pr           上一首                          │
│                                                    │
│ 播放列表:                                          │
│   add/a  <路径>     添加到播放列表                  │
│   list/l            查看播放列表                    │
│   clear/c           清空播放列表                    │
│                                                    │
│ 播放设置:                                          │
│   volume/v <0-255>  设置音量                        │
│   seek/sk <秒>      跳转到指定位置                  │
│   speed/sp <0.5-2>  设置播放速度                    │
│   subtitle/sub <文件> 加载字幕                      │
│                                                    │
│ 播放器管理:                                        │
│   player list        列出播放器                     │
│   player switch <名> 切换播放器                     │
│   player detect      检测可用播放器                 │
│                                                    │
│ 系统:                                              │
│   status/st         查看播放状态                    │
│   cache [clear]     URL缓存管理                     │
│   help/h/?          显示此帮助                      │
│   exit/e       	  退出控制器                      │
│   quit/q/exit       退出控制器                      │
└────────────────────────────────────────────────────┘`)
}

