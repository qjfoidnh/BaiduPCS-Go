package transfer

import (
	"time"
)

type (
	//DownloadInstanceInfo 状态详细信息, 用于导出状态文件
	DownloadInstanceInfo struct {
		DownloadStatus *DownloadStatus
		Ranges         RangeList
	}

	// DownloadInstanceInfoExporter 断点续传类型接口
	DownloadInstanceInfoExporter interface {
		GetInstanceInfo() *DownloadInstanceInfo
		SetInstanceInfo(*DownloadInstanceInfo)
	}
)

// GetInstanceInfo 从断点信息获取下载状态
func (m *DownloadInstanceInfoExport) GetInstanceInfo() (eii *DownloadInstanceInfo) {
	eii = &DownloadInstanceInfo{
		Ranges: m.Ranges,
	}

	var downloaded int64
	switch m.RangeGenMode {
	case RangeGenMode_BlockSize:
		downloaded = m.GenBegin - eii.Ranges.Len()
	default:
		downloaded = m.TotalSize - eii.Ranges.Len()
	}
	eii.DownloadStatus = &DownloadStatus{
		startTime:  time.Now(),
		totalSize:  m.TotalSize,
		downloaded: downloaded,
		gen:        NewRangeListGenBlockSize(m.TotalSize, m.GenBegin, m.BlockSize),
	}
	switch m.RangeGenMode {
	case RangeGenMode_BlockSize:
		eii.DownloadStatus.gen = NewRangeListGenBlockSize(m.TotalSize, m.GenBegin, m.BlockSize)
	default:
		eii.DownloadStatus.gen = NewRangeListGenDefault(m.TotalSize, m.TotalSize, len(m.Ranges), len(m.Ranges))
	}
	return eii
}

// SetInstanceInfo 从下载状态导出断点信息
func (m *DownloadInstanceInfoExport) SetInstanceInfo(eii *DownloadInstanceInfo) {
	if eii == nil {
		return
	}

	if eii.DownloadStatus != nil {
		m.TotalSize = eii.DownloadStatus.TotalSize()
		if eii.DownloadStatus.gen != nil {
			m.GenBegin = eii.DownloadStatus.gen.LoadBegin()
			m.BlockSize = eii.DownloadStatus.gen.LoadBlockSize()
			m.RangeGenMode = eii.DownloadStatus.gen.RangeGenMode()
		} else {
			m.RangeGenMode = RangeGenMode_Default
		}
	}
	m.Ranges = eii.Ranges
}
