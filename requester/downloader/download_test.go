package downloader

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsverbose"
	"github.com/iikira/BaiduPCS-Go/requester"
	"os"
	"testing"
	"time"
)

var (
	url1 = "https://dldir1.qq.com/qqfile/qq/TIM2.1.8/23475/TIM2.1.8.exe"
	url2 = "https://git.oschina.net/lufenping/pixabay_img/raw/master/tiny-20170712/lizard-2427248_1920.jpg"
)

func TestRandomNumber(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(RandomNumber(0, 5))
	}
}

func TestExample(t *testing.T) {
	DoDownload(url2, "lizard-2427248_1920.jpg", nil)
}

func TestDownloadTIM(t *testing.T) {
	pcsverbose.IsVerbose = true

	file, _ := os.OpenFile("tim.exe", os.O_CREATE|os.O_WRONLY, 0777)
	d := NewDownloader(url1, file, &Config{
		MaxParallel:       10,
		CacheSize:         8192,
		InstanceStatePath: "tmp.txt",
	})

	client := requester.NewHTTPClient()
	client.SetTimeout(10 * time.Second)
	d.SetClient(client)

	go func() {
		for {
			if d.monitor != nil {
				fmt.Println(d.monitor.ShowWorkers())
			}
			time.Sleep(1e9)
		}
	}()
	go func() {
		time.Sleep(3e9)
		d.Pause()
		time.Sleep(5e9)
		d.Resume()
		time.Sleep(9e9)
		d.Pause()
		time.Sleep(5e9)
		d.Resume()
		time.Sleep(3e9)
		d.Cancel()
		fmt.Println("canceled")
		time.Sleep(3e9)
	}()
	err := d.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func newSlice() [][]byte {
	s := make([][]byte, 20)
	s[0] = []byte("kjashdfiuqwheirhwuq")
	s[9] = []byte("kjashdfiuqwheirhwuq")
	return s
}

func rangeSlice(f func(key int, by []byte) bool) {
	s := newSlice()
	for k := range s {
		if s[k] == nil {
			continue
		}
		if !f(k, s[k]) {
			break
		}
	}
}

func BenchmarkRange1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var a = 0
		rangeSlice(func(key int, s []byte) bool {
			a++
			return true
		})
	}
}

func BenchmarkRange2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := newSlice()
		a := 0
		for k := range s {
			if s[k] == nil {
				continue
			}
			a++
		}
	}
}
