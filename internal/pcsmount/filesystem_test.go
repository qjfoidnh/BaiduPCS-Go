//go:build windows || cgo

package pcsmount

import (
	"testing"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/winfsp/cgofuse/fuse"
)

func TestCleanRemotePath(t *testing.T) {
	tests := map[string]string{
		"":               "/",
		"/":              "/",
		"我的资源":           "/我的资源",
		`\我的资源\视频\..\文档`: "/我的资源/文档",
		"/a/../../b":     "/b",
	}
	for input, want := range tests {
		if got := cleanRemotePath(input); got != want {
			t.Errorf("cleanRemotePath(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestRemotePathStaysUnderMountRoot(t *testing.T) {
	fs := &FileSystem{root: "/我的资源"}
	if got, want := fs.remotePath("/视频/a.mp4"), "/我的资源/视频/a.mp4"; got != want {
		t.Fatalf("remotePath() = %q, want %q", got, want)
	}
	if got, want := fs.remotePath("/../../文档"), "/我的资源/文档"; got != want {
		t.Fatalf("remotePath() traversal result = %q, want %q", got, want)
	}
}

func TestFillStatReadOnlyModes(t *testing.T) {
	fileStat := &fuse.Stat_t{}
	fillStat(&baidupcs.FileDirectory{FsID: 7, Size: 513}, fileStat)
	if fileStat.Mode != fuse.S_IFREG|0444 || fileStat.Blocks != 2 {
		t.Fatalf("unexpected file stat: mode=%o blocks=%d", fileStat.Mode, fileStat.Blocks)
	}

	dirStat := &fuse.Stat_t{}
	fillStat(&baidupcs.FileDirectory{FsID: 8, Isdir: true}, dirStat)
	if dirStat.Mode != fuse.S_IFDIR|0555 || dirStat.Nlink != 2 {
		t.Fatalf("unexpected directory stat: mode=%o nlink=%d", dirStat.Mode, dirStat.Nlink)
	}
}
