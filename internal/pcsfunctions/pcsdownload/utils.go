package pcsdownload

import (
	"encoding/hex"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

// CheckFileValid 检测文件有效性
func CheckFileValid(filePath string, fileInfo *baidupcs.FileDirectory) error {
	if len(fileInfo.BlockList) != 1 {
		return ErrDownloadNotSupportChecksum
	}

	f := checksum.NewLocalFileChecksum(filePath, int(baidupcs.SliceMD5Size))
	err := f.OpenPath()
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Sum(checksum.CHECKSUM_MD5)
	if err != nil {
		return err
	}
	md5Str := hex.EncodeToString(f.MD5)

	if md5Str != fileInfo.MD5 { // md5不一致
		// 检测是否为违规文件
		if IsSkipMd5Checksum(f.Length, md5Str) {
			return ErrDownloadFileBanned
		}
		return ErrDownloadChecksumFailed
	}
	return nil
}

// FileExist 检查文件是否存在,
// 只有当文件存在, 文件大小不为0或断点续传文件不存在时, 才判断为存在
func FileExist(path string) bool {
	if info, err := os.Stat(path); err == nil {
		if info.Size() == 0 {
			return false
		}
		if _, err = os.Stat(path + DownloadSuffix); err != nil {
			return true
		}
	}

	return false
}

// FixHTTPLinkURL 通过配置, 确定链接使用的协议(http,https)
func FixHTTPLinkURL(linkURL *url.URL) {
	if pcsconfig.Config.EnableHTTPS {
		if linkURL.Scheme == "http" {
			linkURL.Scheme = "https"
		}
	}
}

func CloneJarWithDomain(srcJar http.CookieJar, newURL string) (http.CookieJar, error) {
	if srcJar == nil {
		return nil, fmt.Errorf("srcJar is nil")
	}
	dstJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(newURL)
	newDomain := u.Hostname()
	u, _ = url.Parse("https://" + pcsconfig.Config.PCSAddr + "/")
	cookies := srcJar.Cookies(u)
	for _, c := range cookies {
		nc := *c
		nc.Domain = newDomain
		newURL, _ := url.Parse("https://" + newDomain + "/")
		dstJar.SetCookies(newURL, []*http.Cookie{&nc})
	}
	return dstJar, nil
}
