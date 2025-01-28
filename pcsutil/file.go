package pcsutil

import (
	"github.com/kardianos/osext"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsverbose"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func IsPipeInput() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe
}

// IsIPhoneOS 是否为苹果移动设备
func IsIPhoneOS() bool {
	if runtime.GOOS == "darwin" && (runtime.GOARCH == "arm" || runtime.GOARCH == "arm64") {
		_, err := os.Stat("Info.plist")
		return err == nil
	}
	return false
}

// ChWorkDir 切换回工作目录
func ChWorkDir() {
	if !IsIPhoneOS() {
		return
	}

	dir, err := filepath.Abs("")
	if err != nil {
		return
	}

	subPath := filepath.Dir(os.Args[0])
	os.Chdir(strings.TrimSuffix(dir, subPath))
}

// Executable 获取程序所在的真实目录或真实相对路径
func Executable() string {
	executablePath, err := osext.Executable()
	if err != nil {
		pcsverbose.Verbosef("DEBUG: osext.Executable: %s\n", err)
		executablePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			pcsverbose.Verbosef("DEBUG: filepath.Abs: %s\n", err)
			executablePath = filepath.Dir(os.Args[0])
		}
	}

	if IsIPhoneOS() {
		executablePath = filepath.Join(strings.TrimSuffix(executablePath, os.Args[0]), filepath.Base(os.Args[0]))
	}

	// 读取链接
	linkedExecutablePath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		pcsverbose.Verbosef("DEBUG: filepath.EvalSymlinks: %s\n", err)
		return executablePath
	}
	return linkedExecutablePath
}

// ExecutablePath 获取程序所在目录
func ExecutablePath() string {
	return filepath.Dir(Executable())
}

// ExecutablePathJoin 返回程序所在目录的子目录
func ExecutablePathJoin(subPath string) string {
	return filepath.Join(ExecutablePath(), subPath)
}

// WalkDir 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
// 支持 Linux/macOS 软链接
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 32)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	var walkFunc fs.WalkDirFunc
	walkFunc = func(filename string, fi fs.DirEntry, err error) error { //遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() { // 忽略目录和空文件
			return nil
		}
		fileInfo, err := fi.Info()
		if err != nil || fileInfo.Size() == 0 {
			return nil
		}
		if fileInfo.Mode()&os.ModeSymlink != 0 { // 读取 symbol link
			targetFileInfo, _ := os.Stat(filename)
			if targetFileInfo.IsDir() {
				err = filepath.WalkDir(filename+string(os.PathSeparator), walkFunc)
				return err
			}
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, path.Clean(filename))
		}
		return nil
	}

	err = filepath.WalkDir(dirPth, walkFunc)
	return files, err
}

// ConvertToUnixPathSeparator 将 windows 目录分隔符转换为 Unix 的
func ConvertToUnixPathSeparator(p string) string {
	return strings.Replace(p, "\\", "/", -1)
}

func ChPathLegal(p string) bool {
	illegal_chars := "<>|:\"*?,\\"
	if runtime.GOOS == "windows" {
		illegal_chars = "<>|\"*?,\\"
	}

	if strings.ContainsAny(p, illegal_chars) {
		return false
	}
	return true
}
