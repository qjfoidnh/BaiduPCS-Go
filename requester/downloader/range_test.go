package downloader_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/requester/downloader"
	"testing"
)

func TestRangeListGen(t *testing.T) {
	gen1 := downloader.NewRangeListGenDefault(1024, 0, 0, 10)
	gen2 := downloader.NewRangeListGenBlockSize(1024, 0, 53)

	for mode, gen := range []*downloader.RangeListGen{gen1, gen2} {
		fmt.Printf("[%d] ----\n", mode+1)
		for i, r := gen.GenRange(); r != nil; i, r = gen.GenRange() {
			fmt.Printf("%d: %s\n", i, r.ShowDetails())
		}
		fmt.Println()
	}
}
