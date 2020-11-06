package pcsupdate

type (
	// AssetInfo asset 信息
	AssetInfo struct {
		Name               string `json:"name"`
		ContentType        string `json:"content_type"`
		State              string `json:"state"`
		Size               int64  `json:"size"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}

	// ReleaseInfo 发布信息
	ReleaseInfo struct {
		TagName string       `json:"tag_name"`
		Assets  []*AssetInfo `json:"assets"`
	}
)
