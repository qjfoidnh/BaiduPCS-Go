package pcscommand

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

// RunShareTransfer 执行分享链接转存到网盘
func RunShareTransfer(params []string, opt *baidupcs.TransferOption) {
	var link = params[0]
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接格式非法")
		return
	}
	queryParams := parsedURL.Query()
	extraCode := queryParams.Get("pwd")
	if len(params) == 1 {
		if strings.Contains(link, "bdlink=") || !strings.Contains(link, "pan.baidu.com/") {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "秒传已不再被支持")
			return
		}
	} else if len(params) == 2 {
		extraCode = params[1]
	}

	featureStr := path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
	if featureStr == "init" {
		featureStr = "1" + queryParams.Get("surl")
	}
	if len(featureStr) > 23 || featureStr[0:1] != "1" || len(extraCode) != 4 {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接地址或提取码非法")
		return
	}
	pcs := GetBaiduPCS()
	tokens := pcs.AccessSharePage(featureStr, true)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}

	verifyUrl := pcs.GenerateShareQueryURL("verify", map[string]string{
		"shareid":    tokens["shareid"],
		"time":       strconv.Itoa(int(time.Now().UnixMilli())),
		"clienttype": "1",
		"uk":         tokens["share_uk"],
	}).String()
	res := pcs.PostShareQuery(verifyUrl, link, map[string]string{
		"pwd":       extraCode,
		"vcode":     "null",
		"vcode_str": "null",
		"bdstoken":  tokens["bdstoken"],
	})
	if res["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"])
		return
	}

	pcs.UpdatePCSCookies(true)

	tokens = pcs.AccessSharePage(featureStr, false)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}
	featureMap := map[string]string{
		"bdstoken": tokens["bdstoken"],
		"root":     "1",
		"web":      "5",
		"app_id":   baidupcs.PanAppID,
		"shorturl": featureStr[1:],
		"channel":  "chunlei",
	}
	queryShareInfoUrl := pcs.GenerateShareQueryURL("list", featureMap).String()
	transMetas := pcs.ExtractShareInfo(queryShareInfoUrl, tokens["shareid"], tokens["share_uk"], tokens["bdstoken"])

	if transMetas["ErrMsg"] != "success" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, transMetas["ErrMsg"])
		return
	}
	transMetas["path"] = GetActiveUser().Workdir
	if transMetas["item_num"] != "1" && opt.Collect {
		transMetas["filename"] += "等文件"
		transMetas["path"] = path.Join(GetActiveUser().Workdir, transMetas["filename"])
		pcs.Mkdir(transMetas["path"])
	}
	transMetas["referer"] = "https://pan.baidu.com/s/" + featureStr
	pcs.UpdatePCSCookies(true)
	resp := pcs.GenerateRequestQuery("POST", transMetas)
	if resp["ErrNo"] != "0" {
		if resp["ErrMsg"] != "" && (strings.Contains(resp["ErrMsg"], "转存文件数") || strings.Contains(resp["ErrMsg"], "单次最大转存数")) {
			pcsCommandVerbose.Infof("调试:: 直接转存失败，错误=%s，开始分批转存\n", resp["ErrMsg"])
			runShareTransferByBatch(pcs, featureStr, tokens["shareid"], tokens["share_uk"], tokens["bdstoken"], transMetas["path"], opt)
			return
		}
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, resp["ErrMsg"])
		return
	}
	if opt.Collect {
		resp["filename"] = transMetas["filename"]
	}
	fmt.Printf("%s成功, 保存了%s到当前目录\n", baidupcs.OperationShareFileSavetoLocal, resp["filename"])
	if opt.Download {
		fmt.Println("10s后开始下载")
		time.Sleep(10 * time.Second)
		paths := strings.Split(resp["filenames"], ",")
		RunDownload(paths, nil)
	}
}

func runShareTransferByBatch(pcs *baidupcs.BaiduPCS, featureStr, shareID, shareUK, bdstoken, targetPath string, opt *baidupcs.TransferOption) {
	shortURL := featureStr[1:]
	shareIDInt, err := strconv.ParseInt(shareID, 10, 64)
	if err != nil {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "分享ID解析失败")
		return
	}
	fileChan, errMsg := pcs.GetShareFileListEx(shareIDInt, shareUK, shortURL, bdstoken, opt, true, true)

	if errMsg != "" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, errMsg)
		return
	}

	pcsCommandVerbose.Infof("fileChan=%d,%v\n", len(fileChan), fileChan)

	pcsCommandVerbose.Infof("调试: 分批转存开始 shareID=%d shortURL=%s path=%s\n", shareIDInt, shortURL, targetPath)

	batchSize := 20
	maxFilesPerBatch := 500
	parallel := opt.Parallel
	if parallel <= 0 {
		parallel = 3
	}

	pcsCommandVerbose.Infof("调试: 分批转存参数 batchSize=%d, maxFilesPerBatch=%d, parallel=%d\n", batchSize, maxFilesPerBatch, parallel)
	pcsCommandVerbose.Infof("每批转存最多 %d 个目录，最多 %d 个文件，并发数: %d\n", batchSize, maxFilesPerBatch, parallel)

	type transferStat struct {
		success int
		failed  int
		mu      sync.Mutex
	}
	stat := &transferStat{}

	var mu sync.Mutex
	processed := 0
	failCount := 0
	skipped := 0
	failedFiles := make([]string, 0, 1000)

	sem := make(chan struct{}, parallel)
	var wg sync.WaitGroup

	shouldSkip := func(dirName string) bool {
		var skipped = opt.SkipRange != "" && strings.Contains(dirName, opt.SkipRange)
		return skipped
	}

	isDuplicateError := func(errMsg string) bool {
		return strings.Contains(errMsg, "当前目录下已有") || strings.Contains(errMsg, "文件重复") || strings.Contains(errMsg, "已存在") || strings.Contains(errMsg, "重复")
	}

	transferSingle := func(file *baidupcs.ShareFileInfo, parentPath string) (bool, string) {
		if shouldSkip(parentPath) {
			pcsCommandVerbose.Infof("调试6: 跳过目录 %s, 匹配 skiprange=%s\n", parentPath, opt.SkipRange)
			return false, "skipped"
		}
		referer := "https://pan.baidu.com/s/" + featureStr
		for attempt := 1; attempt <= 3; attempt++ {
			_, errMsg := pcs.BatchTransferShortURL(shortURL, referer, shareID, shareUK, bdstoken, []int64{file.FsID}, parentPath, file.Filename)
			if errMsg == "" || isDuplicateError(errMsg) {
				return true, ""
			}
			if attempt < 3 {
				time.Sleep(500 * time.Millisecond)
			}
			if attempt == 3 {
				return false, errMsg
			}
		}
		return false, "未知错误"
	}

	transferGroup := func(files []*baidupcs.ShareFileInfo, parentPath string) (int, int, []string) {
		if shouldSkip(parentPath) {
			pcsCommandVerbose.Infof("调试4: 跳过目录 %s, 匹配 skiprange=%s\n", parentPath, opt.SkipRange)
			return 0, 0, nil
		}
		referer := "https://pan.baidu.com/s/" + featureStr
		fsIDs := make([]int64, 0, len(files))
		for _, f := range files {
			fsIDs = append(fsIDs, f.FsID)
		}
		filename := ""
		if len(files) == 1 {
			filename = files[0].Filename
		}

		_, errMsg := pcs.BatchTransferShortURL(shortURL, referer, shareID, shareUK, bdstoken, fsIDs, parentPath, filename)
		if errMsg == "" {
			return len(files), 0, nil
		}

		failedPaths := make([]string, 0, len(files))
		failed := 0
		success := 0
		for _, f := range files {
			if shouldSkip(f.Path) {
				pcsCommandVerbose.Infof("调试1: 跳过目录 %s, 匹配 skiprange=%s\n", f.Path, opt.SkipRange)
				continue

			}

			ok, err := transferSingle(f, parentPath)
			if ok {
				success++
			} else {
				failed++
				if len(failedPaths) < 1000 {
					failedPaths = append(failedPaths, f.Path)
				}
				pcsCommandVerbose.Infof("\n转存单目录失败，重试后仍然失败: %s -> %s\n", f.Path, err)
			}
		}

		return success, failed, failedPaths
	}

	commitBatch := func(batch []*baidupcs.ShareFileInfo) {
		pcsCommandVerbose.Infof("commitBatch: %v\n", batch)
		sem <- struct{}{}
		wg.Add(1)
		go func(files []*baidupcs.ShareFileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()

			groups := make(map[string][]*baidupcs.ShareFileInfo)
			groupOrder := make([]string, 0, len(files))
			for _, f := range files {
				if shouldSkip(f.Path) {
					pcsCommandVerbose.Infof("调试5: 跳过目录 %s, 匹配 skiprange=%s\n", f.Path, opt.SkipRange)

					continue
				}
				relativePath := strings.TrimPrefix(f.Path, "/")
				parent := path.Dir(path.Join(targetPath, relativePath))
				if parent == "." {
					parent = targetPath
				}
				if _, ok := groups[parent]; !ok {
					groupOrder = append(groupOrder, parent)
				}
				groups[parent] = append(groups[parent], f)
			}

			for _, parentPath := range groupOrder {
				if shouldSkip(parentPath) {
					pcsCommandVerbose.Infof("调试3: 跳过目录 %s, 匹配 skiprange=%s\n", parentPath, opt.SkipRange)
					continue
				}

				filesInGroup := groups[parentPath]
				if err := pcs.Mkdir(parentPath); err != nil {
					fmt.Printf("\n创建目标目录失败: %s,%s\n", parentPath, err)
				}

				successCount, failedCount, failedPaths := transferGroup(filesInGroup, parentPath)

				mu.Lock()
				processed += len(filesInGroup)
				failCount += failedCount
				if successCount > 0 {
					stat.mu.Lock()
					stat.success += successCount
					stat.mu.Unlock()
				}
				if len(failedFiles) < 1000 {
					remaining := 1000 - len(failedFiles)
					if remaining > len(failedPaths) {
						remaining = len(failedPaths)
					}
					failedFiles = append(failedFiles, failedPaths[:remaining]...)
				}
				pcsCommandVerbose.Infof("\r转存进度1: 已处理 %d 目录, 成功: %d, 失败: %d", processed, stat.success, failCount)
				mu.Unlock()
			}
		}(batch)
	}

	batch := make([]*baidupcs.ShareFileInfo, 0, batchSize)
	currentFiles := 0

	pcsCommandVerbose.Infof("fileChan2=%d,%v\n", len(fileChan), fileChan)
	pcsCommandVerbose.Infof("batch=%d,%v\n", len(batch), batch)
	for file := range fileChan {
		//dirName := path.Base(file.Path)
		if shouldSkip(file.Path) {
			pcsCommandVerbose.Infof("调试0: 跳过目录 %s, 匹配 skiprange=%s\n", file.Path, opt.SkipRange)
			// mu.Lock()
			// skipped++
			// mu.Unlock()
			// break
			continue
		}
		if len(batch) > 0 && (currentFiles+file.FileCount > maxFilesPerBatch || len(batch) >= batchSize) {
			commitBatch(batch)
			batch = make([]*baidupcs.ShareFileInfo, 0, batchSize)
			currentFiles = 0
		}
		batch = append(batch, file)
		currentFiles += file.FileCount
		if len(batch) >= batchSize {
			commitBatch(batch)
			batch = make([]*baidupcs.ShareFileInfo, 0, batchSize)
			currentFiles = 0
		}
	}

	if len(batch) > 0 {
		commitBatch(batch)
	}

	wg.Wait()

	if processed == 0 {
		fmt.Println("分享链接中没有叶子目录可转存")
		//return
	}

	pcsCommandVerbose.Infof("\n批量转存完成!\n")
	pcsCommandVerbose.Infof("总计: %d 目录, 成功: %d, 失败: %d, 跳过: %d\n", processed, stat.success, failCount, skipped)
	if failCount > 0 {
		fmt.Println("\n失败的目录 (最多显示1000个):")
		for _, f := range failedFiles {
			pcsCommandVerbose.Infof("  - %s\n", f)
		}
		if len(failedFiles) == 1000 && failCount > 1000 {
			pcsCommandVerbose.Infof("  ... 还有 %d 个失败目录\n", failCount-1000)
		}
	}
}

// RunShareTransferBatch 批量转存分享链接中的文件
func RunShareTransferBatch(params []string, opt *baidupcs.TransferOption) {
	var link = params[0]
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接格式非法")
		return
	}
	queryParams := parsedURL.Query()
	extraCode := queryParams.Get("pwd")
	if len(params) == 2 {
		extraCode = params[1]
	}

	featureStr := path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
	if featureStr == "init" {
		featureStr = "1" + queryParams.Get("surl")
	}
	if len(featureStr) > 23 || featureStr[0:1] != "1" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接地址非法")
		return
	}
	shortURL := featureStr[1:]

	pcs := GetBaiduPCS()
	tokens := pcs.AccessSharePage(featureStr, true)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}

	if len(extraCode) == 4 {
		verifyUrl := pcs.GenerateShareQueryURL("verify", map[string]string{
			"shareid":    tokens["shareid"],
			"time":       strconv.Itoa(int(time.Now().UnixMilli())),
			"clienttype": "1",
			"uk":         tokens["share_uk"],
		}).String()
		res := pcs.PostShareQuery(verifyUrl, link, map[string]string{
			"pwd":       extraCode,
			"vcode":     "null",
			"vcode_str": "null",
			"bdstoken":  tokens["bdstoken"],
		})
		if res["ErrMsg"] != "0" {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"])
			return
		}
	}

	pcs.UpdatePCSCookies(true)

	tokens = pcs.AccessSharePage(featureStr, false)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}

	fmt.Println("正在获取分享文件列表...")
	shareID, _ := strconv.ParseInt(tokens["shareid"], 10, 64)
	fileChan, errMsg := pcs.GetShareFileList(shareID, tokens["share_uk"], shortURL, tokens["bdstoken"], opt, true)
	if errMsg != "" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, errMsg)
		return
	}

	targetPath := GetActiveUser().Workdir
	if !opt.NoCollect {
		createDir := path.Join(targetPath, "批量转存_"+time.Now().Format("20060102_150405"))
		fmt.Printf("即将创建目录: %s\n", createDir)
		if err := pcs.Mkdir(createDir); err != nil {
			fmt.Printf("创建目录失败: %s, 将直接转存到工作目录\n", err)
		} else {
			targetPath = createDir
		}
	}

	batchSize := opt.BatchSize
	if batchSize <= 0 {
		batchSize = 500
	}
	parallel := opt.Parallel
	if parallel <= 0 {
		parallel = 3
	}

	pcsCommandVerbose.Infof("每批转存 %d 个文件, 并发数: %d\n", batchSize, parallel)

	type transferStat struct {
		success int
		failed  int
		mu      sync.Mutex
	}
	stat := &transferStat{}

	var mu sync.Mutex
	processed := 0
	failCount := 0
	failedFiles := make([]string, 0, 1000) // limit to 1000

	sem := make(chan struct{}, parallel)
	var wg sync.WaitGroup

	batch := make([]*baidupcs.ShareFileInfo, 0, batchSize)
	for file := range fileChan {
		batch = append(batch, file)
		if len(batch) >= batchSize {
			// process batch
			sem <- struct{}{}
			wg.Add(1)
			go func(b []*baidupcs.ShareFileInfo) {
				defer func() {
					<-sem
					wg.Done()
				}()

				var fsIDs []int64
				var batchFailedFiles []string
				for _, f := range b {
					fsIDs = append(fsIDs, f.FsID)
				}

				referer := "https://pan.baidu.com/s/" + featureStr
				filename := ""
				if len(b) > 0 {
					filename = b[0].Filename
					if len(b) > 1 {
						filename += "等文件"
					}
				}
				count, errMsg := pcs.BatchTransferShortURL(shortURL, referer, tokens["shareid"], tokens["share_uk"], tokens["bdstoken"], fsIDs, targetPath, filename)
				if errMsg != "" {
					for _, f := range b {
						batchFailedFiles = append(batchFailedFiles, f.Filename)
					}
				} else {
					stat.mu.Lock()
					stat.success += count
					stat.mu.Unlock()
				}

				mu.Lock()
				processed += len(b)
				if errMsg != "" {
					failCount += len(b)
					if len(failedFiles) < 1000 {
						failedFiles = append(failedFiles, batchFailedFiles...)
					}
					pcsCommandVerbose.Infof("\r转存进度2: 已处理 %d, 成功: %d, 失败: %d, 错误: %s", processed, stat.success, failCount, errMsg)
				} else {
					pcsCommandVerbose.Infof("\r转存进度3: 已处理 %d, 成功: %d, 失败: %d", processed, stat.success, failCount)
				}
				mu.Unlock()
			}(batch)
			batch = make([]*baidupcs.ShareFileInfo, 0, batchSize)
		}
	}
	// process remaining batch
	if len(batch) > 0 {
		sem <- struct{}{}
		wg.Add(1)
		go func(b []*baidupcs.ShareFileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()

			var fsIDs []int64
			var batchFailedFiles []string
			for _, f := range b {
				fsIDs = append(fsIDs, f.FsID)
			}

			referer := "https://pan.baidu.com/s/" + featureStr
			filename := ""
			if len(b) > 0 {
				filename = b[0].Filename
				if len(b) > 1 {
					filename += "等文件"
				}
			}
			count, errMsg := pcs.BatchTransferShortURL(shortURL, referer, tokens["shareid"], tokens["share_uk"], tokens["bdstoken"], fsIDs, targetPath, filename)
			if errMsg != "" {
				for _, f := range b {
					batchFailedFiles = append(batchFailedFiles, f.Filename)
				}
			} else {
				stat.mu.Lock()
				stat.success += count
				stat.mu.Unlock()
			}

			mu.Lock()
			processed += len(b)
			if errMsg != "" {
				failCount += len(b)
				if len(failedFiles) < 1000 {
					failedFiles = append(failedFiles, batchFailedFiles...)
				}
				pcsCommandVerbose.Infof("\r转存进度4: 已处理 %d, 成功: %d, 失败: %d, 错误: %s", processed, stat.success, failCount, errMsg)
			} else {
				pcsCommandVerbose.Infof("\r转存进度5: 已处理 %d, 成功: %d, 失败: %d", processed, stat.success, failCount)
			}
			mu.Unlock()
		}(batch)
	}

	wg.Wait()
	if processed == 0 {
		pcsCommandVerbose.Info("分享链接中没有文件")
		return
	}
	pcsCommandVerbose.Infof("\n批量转存完成!\n")
	pcsCommandVerbose.Infof("总计: %d, 成功: %d, 失败: %d\n", processed, stat.success, failCount)

	if failCount > 0 {
		pcsCommandVerbose.Infof("\n失败的文件 (最多显示1000个):")
		for _, f := range failedFiles {
			pcsCommandVerbose.Infof("  - %s\n", f)
		}
		if len(failedFiles) == 1000 && failCount > 1000 {
			pcsCommandVerbose.Infof("  ... 还有 %d 个失败文件\n", failCount-1000)
		}
	}
}
