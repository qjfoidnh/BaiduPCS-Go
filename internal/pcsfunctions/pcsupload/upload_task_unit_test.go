package pcsupload

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/taskframework"
)

func attachTaskInfo(unit *UploadTaskUnit) {
	executor := taskframework.NewTaskExecutor()
	executor.Append(unit, 0)
}

func TestRunReportsOversizedFileAsFailure(t *testing.T) {
	unit := &UploadTaskUnit{
		LocalFileChecksum: &checksum.LocalFileChecksum{
			LocalFileMeta: checksum.LocalFileMeta{
				Path:   "too-large.mp4",
				Length: baidupcs.MaxUploadSize + 1,
			},
		},
	}
	attachTaskInfo(unit)

	result := unit.Run()

	if result == nil || result.Succeed || result.Err == nil {
		t.Fatalf("expected explicit failure result, got %#v", result)
	}
}

func TestRunReportsUnreadableFileAsFailure(t *testing.T) {
	unit := &UploadTaskUnit{
		LocalFileChecksum: checksum.NewLocalFileChecksum(
			filepath.Join(t.TempDir(), "missing.mp4"),
			int(baidupcs.SliceMD5Size),
		),
	}
	attachTaskInfo(unit)

	result := unit.Run()

	if result == nil || result.Succeed || result.Err == nil {
		t.Fatalf("expected explicit failure result, got %#v", result)
	}
}

func TestJustGoonReturnsPreparedFailure(t *testing.T) {
	prepared := &taskframework.TaskUnitRunResult{
		ResultMessage: "目标文件大小超过剩余空间",
		Err:           errors.New("insufficient space"),
	}
	unit := &UploadTaskUnit{
		LocalFileChecksum: &checksum.LocalFileChecksum{},
		Step:              JustGoon,
		prepareResult:     prepared,
	}
	attachTaskInfo(unit)

	result := unit.runPreparedStep()

	if result != prepared {
		t.Fatalf("expected prepared result, got %#v", result)
	}
}
