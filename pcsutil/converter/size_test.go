package converter_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"testing"
)

func TestParseFileSizeStr(t *testing.T) {
	for _, v := range []string{"1k", "3.86mb", "4.001Gb", "32"} {
		size, err := converter.ParseFileSizeStr(v)
		if err != nil {
			t.Fatalf("%s\n", err)
		}
		fmt.Println(v, size)
	}
}
