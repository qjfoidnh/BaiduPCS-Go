package checksum_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsutil/checksum"
	"testing"
)

var (
	flagList = []int{
		checksum.CHECKSUM_MD5 | 000000000000000000000000000 | 00000000000000000000000,
		000000000000000000000 | checksum.CHECKSUM_SLICE_MD5 | 00000000000000000000000,
		000000000000000000000 | 000000000000000000000000000 | checksum.CHECKSUM_CRC32,
		checksum.CHECKSUM_MD5 | checksum.CHECKSUM_SLICE_MD5 | 00000000000000000000000,
		000000000000000000000 | checksum.CHECKSUM_SLICE_MD5 | checksum.CHECKSUM_CRC32,
		checksum.CHECKSUM_MD5 | 000000000000000000000000000 | checksum.CHECKSUM_CRC32,
		checksum.CHECKSUM_MD5 | checksum.CHECKSUM_SLICE_MD5 | checksum.CHECKSUM_CRC32,
	}
)

func printFileMeta(meta *checksum.LocalFileMeta) {
	fmt.Printf("slicemd5: %x, md5: %x, crc32: %x %d\n", meta.SliceMD5, meta.MD5, meta.CRC32, meta.CRC32)
}

func TestChecksum(t *testing.T) {
	fmt.Println("--- file.go")
	for _, flag := range flagList {
		lf, err := checksum.GetFileSum("file.go", flag)
		if err != nil {
			t.Fatal(err)
		}
		printFileMeta(&lf.LocalFileMeta)
	}

	fmt.Println("--- /Users/syy/go/src/github.com/iikira/BaiduPCS-Go/BaiduPCS-Go")
	for _, flag := range flagList {
		lf := checksum.NewLocalFileChecksumWithBufSize("/Users/syy/go/src/github.com/iikira/BaiduPCS-Go/BaiduPCS-Go", checksum.DefaultBufSize-3, checksum.DefaultBufSize)
		err := lf.OpenPath()
		if err != nil {
			t.Fatal(err)
		}
		lf.Sum(flag)
		printFileMeta(&lf.LocalFileMeta)
	}
}
