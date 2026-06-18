package pcscommand

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsdownload"
	"github.com/urfave/cli"
)

const (
	cacheExpireHours = 24
)

// 支持的音频格式
var audioExtensions = map[string]bool{
	"mp3": true, "flac": true, "wav": true, "aac": true,
	"ogg": true, "m4a": true, "wma": true, "ape": true,
	"opus": true, "wv": true, "tta": true,
}

// 已知的文件扩展名（用于判断是文件还是目录）
var knownExtensions = map[string]bool{
	// 音频
	"mp3": true, "flac": true, "wav": true, "aac": true,
	"ogg": true, "m4a": true, "wma": true, "ape": true,
	"opus": true, "wv": true, "tta": true, "aiff": true,
	// 视频
	"mp4": true, "mkv": true, "avi": true, "mov": true,
	"flv": true, "wmv": true, "webm": true, "m4v": true,
	"3gp": true, "rmvb": true, "ts": true,
	// 文档
	"txt": true, "pdf": true, "doc": true, "docx": true,
	"xls": true, "xlsx": true, "ppt": true, "pptx": true,
	"zip": true, "rar": true, "7z": true, "tar": true, "gz": true,
	"jpg": true, "jpeg": true, "png": true, "gif": true, "bmp": true,
	"apk": true, "exe": true, "iso": true,
}

// isFileByExtension 通过扩展名判断是否是文件
func isFileByExtension(path string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	return knownExtensions[ext]
}

// CacheItem 缓存项
type CacheItem struct {
	URL      string
	FilePath string
	ExpireAt time.Time
}

// URLCache URL 缓存
type URLCache struct {
	mu    sync.RWMutex
	files map[string]CacheItem
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
		ExpireAt: time.Now().Add(time.Duration(cacheExpireHours) * time.Hour),
	}
}

func (c *URLCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files = make(map[string]CacheItem)
}

// Play 播放入口
func Play(c *cli.Context) error {
	if c.NArg() == 0 {
		return controller()
	}

	path := c.Args().Get(0)
	return playPath(path)
}

// PlayFile 播放文件或目录
func PlayFile(path string) error {
	return playPath(path)
}

// playPath 智能播放：有扩展名=文件，无扩展名=尝试目录
func playPath(path string) error {
	// 确保路径以 / 开头
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 有已知扩展名 -> 文件播放
	if isFileByExtension(path) {
		return playSingleFile(path)
	}

	// 无扩展名或以 / 结尾 -> 尝试目录播放
	if !strings.Contains(filepath.Base(path), ".") || strings.HasSuffix(path, "/") {
		err := playDirectory(path)
		if err == nil {
			return nil
		}
		// 目录播放失败，尝试作为文件
		return playSingleFile(path)
	}

	// 默认尝试文件播放
	return playSingleFile(path)
}

// playSingleFile 播放单个文件
func playSingleFile(path string) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if !audioExtensions[ext] {
		return fmt.Errorf("不支持的音频格式: .%s", ext)
	}

	fmt.Printf("正在播放: %s\n", filepath.Base(path))

	pcs := GetBaiduPCS()
	url, err := getURL(pcs, path)
	if err != nil {
		return fmt.Errorf("获取URL失败: %v", err)
	}

	return vlcCommand("add " + url)
}

// playDirectory 播放目录下所有音频文件
func playDirectory(dirPath string) error {
	// 确保路径以 / 开头
	if !strings.HasPrefix(dirPath, "/") {
		dirPath = "/" + dirPath
	}
	// 去掉末尾的 /
	dirPath = strings.TrimRight(dirPath, "/")

	fmt.Printf("正在扫描目录: %s\n", dirPath)

	// 获取目录下的文件列表
	pcs := GetBaiduPCS()
	fileList, err := pcs.FilesDirectoriesList(dirPath, nil)
	if err != nil {
		return fmt.Errorf("获取目录列表失败: %v", err)
	}

	// 筛选音频文件
	var audioFiles []string
	for _, file := range fileList {
		if file.Isdir {
			continue
		}
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
		if audioExtensions[ext] {
			audioFiles = append(audioFiles, file.Filename)
		}
	}

	if len(audioFiles) == 0 {
		return fmt.Errorf("目录下没有支持的音频文件")
	}

	// 按文件名排序
	sort.Strings(audioFiles)

	fmt.Printf("找到 %d 个音频文件:\n", len(audioFiles))
	for i, f := range audioFiles {
		fmt.Printf("  %d. %s\n", i+1, f)
	}

	// 清空当前播放列表
	vlcCommand("clear")

	fmt.Println("\n正在添加到播放列表...")
	successCount := 0

	// 初始化播放列表
    currentPlaylist = &PlaylistData{Files: make([]PlaylistFile, 0)}

	for _, fileName := range audioFiles {
		fullPath := dirPath + "/" + fileName
		url, err := getURL(pcs, fullPath)
		if err != nil {
			fmt.Printf("  跳过 %s: %v\n", fileName, err)
			continue
		}

		// 保存到播放列表数据
        currentPlaylist.Files = append(currentPlaylist.Files, PlaylistFile{
            Path: fullPath, URL: url, Name: fileName,
        })
		
		if successCount == 0 {
			// 第一个文件直接播放
			vlcCommand("add " + url)
			fmt.Printf("  ▶ 开始播放: %s\n", fileName)
		} else {
			// 后续文件加入队列
			if err := vlcCommand("enqueue " + url); err != nil {
				fmt.Printf("  添加失败 %s: %v\n", fileName, err)
				continue
			}
		}
		successCount++
	}

	SavePlaylist() 

	if successCount == 0 {
		return fmt.Errorf("没有文件可以播放")
	}

	fmt.Printf("\n成功添加 %d 个文件到播放列表\n", successCount)
	fmt.Println("开始播放...")

	//vlcCommand("play")

	return nil
}

// getURL 获取文件下载链接（带缓存）
func getURL(pcs *baidupcs.BaiduPCS, path string) (string, error) {
	// 确保路径以 / 开头
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 检查缓存
	if url, ok := urlCache.Get(path); ok {
		return url, nil
	}

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
	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║     BaiduPCS-Go 多媒体播放器             ║")
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Println("║  输入 help 查看命令                      ║")
	fmt.Println("║  输入 quit 退出控制器                    ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	/*
	// 恢复上次播放列表
    if pl := LoadPlaylist(); pl != nil && len(pl.Files) > 0 {
        fmt.Printf("发现上次播放列表 (%d 首)，正在恢复...\n", len(pl.Files))
        currentPlaylist = pl
        pcs := GetBaiduPCS()
        vlcCommand("clear")
        for i, f := range pl.Files {
            url, _ := getURL(pcs, f.Path)
            if url == "" {
                url = f.URL
            }
            if i == 0 {
                vlcCommand("add " + url)
            } else {
                vlcCommand("enqueue " + url)
            }
        }
    }
	*/

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
			vlcCommand("pause")
			fmt.Println("暂停/继续")
		case "stop", "s":
			vlcCommand("stop")
			fmt.Println("已停止")
		case "next", "n":
			vlcCommand("next")
			fmt.Println("下一首")
		case "prev", "pr":
			vlcCommand("prev")
			fmt.Println("上一首")
		case "list", "l":
			vlcCommand("playlist")
		case "clear", "c":
			vlcCommand("clear")
			fmt.Println("播放列表已清空")
		case "status", "st":
			vlcCommand("status")
		case "volume", "v":
			handleVolume(parts)
		case "seek", "sk":
			handleSeek(parts)
		case "speed", "sp":
			handleSpeed(parts)
		case "random", "rd":
			vlcCommand("random on")
			fmt.Println("随机播放: 开")
		case "repeat", "rp":
			handleRepeat(parts)
		case "help", "h", "?":
			showHelp()
		case "quit", "q", "exit":
			fmt.Println("退出播放器控制器")
			return nil
		default:
			fmt.Printf("未知命令: %s (输入 help 查看帮助)\n", cmd)
		}

		fmt.Print("player:> ")
	}

	return scanner.Err()
}

func handleAdd(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: add <路径>")
		fmt.Println("示例: add /音乐/歌曲.mp3")
		fmt.Println("      add /音乐/华语经典")
		return
	}

	path := strings.Join(parts[1:], " ")

	// 确保绝对路径
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 有扩展名 -> 单文件
	if isFileByExtension(path) {
		pcs := GetBaiduPCS()
		url, err := getURL(pcs, path)
		if err != nil {
			fmt.Printf("添加失败: %v\n", err)
			return
		}
		if err := vlcCommand("enqueue " + url); err != nil {
			fmt.Printf("添加失败: %v\n", err)
		} else {
			fmt.Printf("已添加: %s\n", filepath.Base(path))
		}
		return
	}

	// 无扩展名 -> 目录
	path = strings.TrimRight(path, "/")
	pcs := GetBaiduPCS()
	fileList, err := pcs.FilesDirectoriesList(path, nil)
	if err != nil {
		fmt.Printf("获取目录失败: %v\n", err)
		return
	}

	count := 0
	for _, file := range fileList {
		if file.Isdir {
			continue
		}
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
		if audioExtensions[ext] {
			fullPath := path + "/" + file.Filename
			url, err := getURL(pcs, fullPath)
			if err != nil {
				continue
			}
			vlcCommand("enqueue " + url)
			count++
		}
	}
	fmt.Printf("已添加 %d 个文件\n", count)
}

func handlePlay(parts []string) {
	if len(parts) < 2 {
		vlcCommand("play")
		fmt.Println("继续播放")
		return
	}

	path := strings.Join(parts[1:], " ")
	playPath(path)
}

func handleVolume(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: volume <0-255>")
		return
	}
	var volume int
	_, err := fmt.Sscanf(parts[1], "%d", &volume)
	if err != nil || volume < 0 || volume > 255 {
		fmt.Println("音量范围: 0-255")
		return
	}
	vlcCommand(fmt.Sprintf("volume %d", volume))
	fmt.Printf("音量: %d\n", volume)
}

func handleSeek(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: seek <秒>")
		return
	}
	vlcCommand(fmt.Sprintf("seek %s", parts[1]))
}

func handleSpeed(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: speed <0.5-2.0>")
		return
	}
	vlcCommand(fmt.Sprintf("rate %s", parts[1]))
}

func handleRepeat(parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: repeat [off|one|all]")
		return
	}
	switch parts[1] {
	case "off":
		vlcCommand("repeat off")
	case "one":
		vlcCommand("repeat one")
	case "all":
		vlcCommand("repeat all")
	default:
		fmt.Println("用法: repeat [off|one|all]")
	}
}

func showHelp() {
	fmt.Println(`
┌────────────────────────────────────────────────────┐
│                  可用命令列表                       │
├────────────────────────────────────────────────────┤
│ 播放控制:                                          │
│   play/p  <路径>    播放文件或目录                  │
│   pause/pa          暂停/继续                       │
│   stop/s            停止播放                        │
│   next/n            下一首                          │
│   prev/pr           上一首                          │
│                                                    │
│ 播放列表:                                          │
│   add/a  <路径>     添加文件/目录到播放列表          │
│   list/l            查看播放列表                    │
│   clear/c           清空播放列表                    │
│   random/rd         切换随机播放                    │
│   repeat/rp [模式]  重复模式(off/one/all)           │
│                                                    │
│ 播放设置:                                          │
│   volume/v <0-255>  设置音量                        │
│   seek/sk <秒>      跳转到指定位置                  │
│   speed/sp <0.5-2>  设置播放速度                    │
│                                                    │
│ 其他:                                              │
│   status/st         查看播放状态                    │
│   help/h/?          显示此帮助                      │
│   quit/q/exit       退出控制器                      │
│                                                    │
│ 智能识别:                                          │
│   有扩展名(.mp3等) = 文件                          │
│   无扩展名 = 目录                                  │
│                                                    │
│ 示例:                                              │
│   play /音乐/歌曲.mp3      播放单文件               │
│   play /音乐/华语经典      播放整个目录             │
│   add /音乐/              添加整个目录到播放列表    │
└────────────────────────────────────────────────────┘`)
}

// VLC 控制
var (
	vlcCmd   *exec.Cmd
	vlcStdin io.WriteCloser
	vlcStdout io.ReadCloser 
)

type PlaylistFile struct {
    Path string `json:"path"`
    URL  string `json:"url"`
    Name string `json:"name"`
}

type PlaylistData struct {
    Files   []PlaylistFile `json:"files"`
    Current int            `json:"current"`
}

var (
    currentPlaylist *PlaylistData
    playlistFile    = filepath.Join(os.TempDir(), "baidupcs_playlist.json")
)

func startVlc() error {
	if vlcCmd != nil {
		return nil
	}

	cmd := exec.Command("vlc", "--intf", "rc", "--rc-fake-tty", "--quiet")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建标准输入管道失败: %v", err)
	}

	stdout, err := cmd.StdoutPipe() 
    if err != nil {
        return fmt.Errorf("创建标准输出管道失败: %v", err)
    }

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 VLC 失败: %v", err)
	}

	vlcCmd = cmd
	vlcStdin = stdin
	vlcStdout = stdout

	time.Sleep(2 * time.Second)

	fmt.Println("VLC 已启动")
	return nil
}

func vlcCommand(args string) error {
	if len(args) == 0 {
		return fmt.Errorf("命令参数不能为空")
	}

	if vlcCmd == nil {
		fmt.Printf("正在启动 VLC...\n")
		if err := startVlc(); err != nil {
			return fmt.Errorf("启动 VLC 失败: %v", err)
		}
	}

	if vlcStdin == nil {
		return fmt.Errorf("VLC 标准输入管道不可用")
	}

	_, err := vlcStdin.Write([]byte(args + "\n"))
	if err != nil {
		return fmt.Errorf("发送命令失败: %v", err)
	}

	return nil
}

func SavePlaylist() {
    if currentPlaylist == nil {
        return
    }
    data, _ := json.Marshal(currentPlaylist)
    os.WriteFile(playlistFile, data, 0644)
}

func LoadPlaylist() *PlaylistData {
    data, err := os.ReadFile(playlistFile)
    if err != nil {
        return nil
    }
    var pl PlaylistData
    if json.Unmarshal(data, &pl) == nil {
        return &pl
    }
    return nil
}



// InitPlayer 初始化播放器
func InitPlayer() {
	// VLC 会在首次使用时自动启动
}
