//go:build windows || cgo

// Package pcsmount exposes a Baidu PCS directory through FUSE.
package pcsmount

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsdownload"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsverbose"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	defaultCacheTTL = 30 * time.Second
	linkCacheTTL    = 5 * time.Minute
	blockSize       = 4096
)

var mountVerbose = pcsverbose.New("PCSMOUNT")

type metadataCacheEntry struct {
	value     *baidupcs.FileDirectory
	expiresAt time.Time
}

type directoryCacheEntry struct {
	value     baidupcs.FileDirectoryList
	expiresAt time.Time
}

type linkCacheEntry struct {
	value     string
	expiresAt time.Time
}

// FileSystem is a read-only view of a Baidu PCS directory.
type FileSystem struct {
	fuse.FileSystemBase

	pcs      *baidupcs.BaiduPCS
	root     string
	cacheTTL time.Duration

	mu          sync.RWMutex
	metadata    map[string]metadataCacheEntry
	directories map[string]directoryCacheEntry
	links       map[string]linkCacheEntry
	quota       int64
	used        int64
	quotaExpiry time.Time
}

// NewFileSystem validates remoteRoot and creates a read-only file system.
func NewFileSystem(pcs *baidupcs.BaiduPCS, remoteRoot string, cacheTTL time.Duration) (*FileSystem, error) {
	if pcs == nil {
		return nil, errors.New("BaiduPCS client is nil")
	}
	remoteRoot = cleanRemotePath(remoteRoot)
	if cacheTTL <= 0 {
		cacheTTL = defaultCacheTTL
	}

	rootInfo, pcsErr := pcs.FilesDirectoriesMeta(remoteRoot)
	if pcsErr != nil {
		return nil, fmt.Errorf("读取挂载目录元信息失败: %w", pcsErr)
	}
	if !rootInfo.Isdir {
		return nil, fmt.Errorf("挂载源不是目录: %s", remoteRoot)
	}

	fs := &FileSystem{
		pcs:         pcs,
		root:        remoteRoot,
		cacheTTL:    cacheTTL,
		metadata:    make(map[string]metadataCacheEntry),
		directories: make(map[string]directoryCacheEntry),
		links:       make(map[string]linkCacheEntry),
	}
	fs.metadata[remoteRoot] = metadataCacheEntry{value: rootInfo, expiresAt: time.Now().Add(cacheTTL)}
	return fs, nil
}

func cleanRemotePath(value string) string {
	value = strings.ReplaceAll(value, "\\", "/")
	if value == "" {
		return "/"
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return path.Clean(value)
}

func (fs *FileSystem) remotePath(fusePath string) string {
	if fusePath == "" || fusePath == "/" {
		return fs.root
	}
	return path.Join(fs.root, strings.TrimPrefix(cleanRemotePath(fusePath), "/"))
}

func (fs *FileSystem) getMetadata(fusePath string) (*baidupcs.FileDirectory, error) {
	remotePath := fs.remotePath(fusePath)
	now := time.Now()
	fs.mu.RLock()
	entry, ok := fs.metadata[remotePath]
	fs.mu.RUnlock()
	if ok && now.Before(entry.expiresAt) {
		return entry.value, nil
	}

	value, pcsErr := fs.pcs.FilesDirectoriesMeta(remotePath)
	if pcsErr != nil {
		return nil, pcsErr
	}
	fs.mu.Lock()
	fs.metadata[remotePath] = metadataCacheEntry{value: value, expiresAt: now.Add(fs.cacheTTL)}
	fs.mu.Unlock()
	return value, nil
}

func (fs *FileSystem) getDirectory(fusePath string) (baidupcs.FileDirectoryList, error) {
	remotePath := fs.remotePath(fusePath)
	now := time.Now()
	fs.mu.RLock()
	entry, ok := fs.directories[remotePath]
	fs.mu.RUnlock()
	if ok && now.Before(entry.expiresAt) {
		return entry.value, nil
	}

	value, pcsErr := fs.pcs.FilesDirectoriesList(remotePath, baidupcs.DefaultOrderOptions)
	if pcsErr != nil {
		return nil, pcsErr
	}
	fs.mu.Lock()
	fs.directories[remotePath] = directoryCacheEntry{value: value, expiresAt: now.Add(fs.cacheTTL)}
	for _, child := range value {
		fs.metadata[child.Path] = metadataCacheEntry{value: child, expiresAt: now.Add(fs.cacheTTL)}
	}
	fs.mu.Unlock()
	return value, nil
}

func fillStat(info *baidupcs.FileDirectory, stat *fuse.Stat_t) {
	stat.Ino = uint64(info.FsID)
	stat.Nlink = 1
	stat.Size = info.Size
	stat.Blksize = blockSize
	stat.Blocks = (info.Size + 511) / 512
	stat.Atim = fuse.NewTimespec(time.Unix(info.Mtime, 0))
	stat.Mtim = fuse.NewTimespec(time.Unix(info.Mtime, 0))
	stat.Ctim = fuse.NewTimespec(time.Unix(info.Ctime, 0))
	stat.Birthtim = stat.Ctim
	if info.Isdir {
		stat.Mode = fuse.S_IFDIR | 0555
		stat.Nlink = 2
	} else {
		stat.Mode = fuse.S_IFREG | 0444
	}
}

func (fs *FileSystem) Getattr(fusePath string, stat *fuse.Stat_t, _ uint64) int {
	info, err := fs.getMetadata(fusePath)
	if err != nil {
		mountVerbose.Warnf("getattr %s: %s\n", fusePath, err)
		return -fuse.ENOENT
	}
	fillStat(info, stat)
	return 0
}

func (fs *FileSystem) Access(fusePath string, mask uint32) int {
	if mask&(fuse.W_OK|fuse.DELETE_OK) != 0 {
		return -fuse.EROFS
	}
	if _, err := fs.getMetadata(fusePath); err != nil {
		return -fuse.ENOENT
	}
	return 0
}

func (fs *FileSystem) Open(fusePath string, flags int) (int, uint64) {
	if flags&fuse.O_ACCMODE != fuse.O_RDONLY || flags&fuse.O_TRUNC != 0 {
		return -fuse.EROFS, ^uint64(0)
	}
	info, err := fs.getMetadata(fusePath)
	if err != nil {
		return -fuse.ENOENT, ^uint64(0)
	}
	if info.Isdir {
		return -fuse.EISDIR, ^uint64(0)
	}
	return 0, uint64(info.FsID)
}

func (fs *FileSystem) Opendir(fusePath string) (int, uint64) {
	info, err := fs.getMetadata(fusePath)
	if err != nil {
		return -fuse.ENOENT, ^uint64(0)
	}
	if !info.Isdir {
		return -fuse.ENOTDIR, ^uint64(0)
	}
	return 0, uint64(info.FsID)
}

func (fs *FileSystem) Readdir(fusePath string, fill func(string, *fuse.Stat_t, int64) bool, ofst int64, _ uint64) int {
	children, err := fs.getDirectory(fusePath)
	if err != nil {
		mountVerbose.Warnf("readdir %s: %s\n", fusePath, err)
		return -fuse.EIO
	}

	type item struct {
		name string
		info *baidupcs.FileDirectory
	}
	items := make([]item, 0, len(children)+2)
	items = append(items, item{name: "."}, item{name: ".."})
	for _, child := range children {
		items = append(items, item{name: child.Filename, info: child})
	}
	if ofst < 0 || ofst > int64(len(items)) {
		return -fuse.EINVAL
	}
	for index := int(ofst); index < len(items); index++ {
		var stat *fuse.Stat_t
		if items[index].info != nil {
			stat = &fuse.Stat_t{}
			fillStat(items[index].info, stat)
		}
		if !fill(items[index].name, stat, int64(index+1)) {
			break
		}
	}
	return 0
}

func (fs *FileSystem) downloadLink(remotePath string, refresh bool) (string, error) {
	now := time.Now()
	if !refresh {
		fs.mu.RLock()
		entry, ok := fs.links[remotePath]
		fs.mu.RUnlock()
		if ok && now.Before(entry.expiresAt) {
			return entry.value, nil
		}
	}

	links, err := pcsdownload.GetLocateDownloadLinks(fs.pcs, remotePath)
	if err != nil {
		return "", err
	}
	if len(links) == 0 {
		return "", errors.New("百度网盘未返回下载链接")
	}
	pcsdownload.FixHTTPLinkURL(links[0])
	value := links[0].String()
	fs.mu.Lock()
	fs.links[remotePath] = linkCacheEntry{value: value, expiresAt: now.Add(linkCacheTTL)}
	fs.mu.Unlock()
	return value, nil
}

func (fs *FileSystem) readRange(remotePath string, buff []byte, offset int64) (int, error) {
	if len(buff) == 0 {
		return 0, nil
	}
	for attempt := 0; attempt < 2; attempt++ {
		link, err := fs.downloadLink(remotePath, attempt > 0)
		if err != nil {
			return 0, err
		}
		client := pcsconfig.Config.PanHTTPClient()
		client.SetTimeout(2 * time.Minute)
		jar, err := pcsdownload.CloneJarWithDomain(fs.pcs.GetClient().Jar, link)
		if err != nil {
			return 0, err
		}
		client.SetCookiejar(jar)
		end := offset + int64(len(buff)) - 1
		resp, err := client.Req(http.MethodGet, link, nil, map[string]string{
			"Range": fmt.Sprintf("bytes=%d-%d", offset, end),
		})
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			if attempt == 0 {
				continue
			}
			return 0, err
		}
		if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
			return 0, nil
		}
		if resp.StatusCode != http.StatusPartialContent && !(offset == 0 && resp.StatusCode == http.StatusOK) {
			if attempt == 0 && (resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized) {
				continue
			}
			return 0, fmt.Errorf("下载请求返回 %s", resp.Status)
		}
		n, readErr := io.ReadFull(resp.Body, buff)
		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			return n, nil
		}
		return n, readErr
	}
	return 0, errors.New("下载链接刷新后仍不可用")
}

func (fs *FileSystem) Read(fusePath string, buff []byte, offset int64, _ uint64) int {
	if offset < 0 {
		return -fuse.EINVAL
	}
	info, err := fs.getMetadata(fusePath)
	if err != nil {
		return -fuse.ENOENT
	}
	if info.Isdir {
		return -fuse.EISDIR
	}
	if offset >= info.Size {
		return 0
	}
	if remaining := info.Size - offset; int64(len(buff)) > remaining {
		buff = buff[:remaining]
	}
	n, err := fs.readRange(fs.remotePath(fusePath), buff, offset)
	if err != nil {
		mountVerbose.Warnf("read %s at %d: %s\n", fusePath, offset, err)
		return -fuse.EIO
	}
	return n
}

func (fs *FileSystem) Statfs(_ string, stat *fuse.Statfs_t) int {
	now := time.Now()
	fs.mu.RLock()
	quota, used, valid := fs.quota, fs.used, now.Before(fs.quotaExpiry)
	fs.mu.RUnlock()
	if !valid {
		var pcsErr error
		quota, used, pcsErr = fs.pcs.QuotaInfo()
		if pcsErr != nil {
			return -fuse.EIO
		}
		fs.mu.Lock()
		fs.quota, fs.used, fs.quotaExpiry = quota, used, now.Add(fs.cacheTTL)
		fs.mu.Unlock()
	}
	free := quota - used
	if free < 0 {
		free = 0
	}
	stat.Bsize = blockSize
	stat.Frsize = blockSize
	stat.Blocks = uint64(quota / blockSize)
	stat.Bfree = uint64(free / blockSize)
	stat.Bavail = stat.Bfree
	stat.Namemax = 255
	return 0
}

func (fs *FileSystem) Release(string, uint64) int        { return 0 }
func (fs *FileSystem) Releasedir(string, uint64) int     { return 0 }
func (fs *FileSystem) Flush(string, uint64) int          { return 0 }
func (fs *FileSystem) Fsync(string, bool, uint64) int    { return 0 }
func (fs *FileSystem) Fsyncdir(string, bool, uint64) int { return 0 }

func (fs *FileSystem) Mknod(string, uint32, uint64) int         { return -fuse.EROFS }
func (fs *FileSystem) Mkdir(string, uint32) int                 { return -fuse.EROFS }
func (fs *FileSystem) Unlink(string) int                        { return -fuse.EROFS }
func (fs *FileSystem) Rmdir(string) int                         { return -fuse.EROFS }
func (fs *FileSystem) Link(string, string) int                  { return -fuse.EROFS }
func (fs *FileSystem) Symlink(string, string) int               { return -fuse.EROFS }
func (fs *FileSystem) Rename(string, string) int                { return -fuse.EROFS }
func (fs *FileSystem) Chmod(string, uint32) int                 { return -fuse.EROFS }
func (fs *FileSystem) Chown(string, uint32, uint32) int         { return -fuse.EROFS }
func (fs *FileSystem) Utimens(string, []fuse.Timespec) int      { return -fuse.EROFS }
func (fs *FileSystem) Create(string, int, uint32) (int, uint64) { return -fuse.EROFS, ^uint64(0) }
func (fs *FileSystem) Truncate(string, int64, uint64) int       { return -fuse.EROFS }
func (fs *FileSystem) Write(string, []byte, int64, uint64) int  { return -fuse.EROFS }
func (fs *FileSystem) Setxattr(string, string, []byte, int) int { return -fuse.EROFS }
func (fs *FileSystem) Removexattr(string, string) int           { return -fuse.EROFS }
