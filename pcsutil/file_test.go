package pcsutil_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"testing"
)

func TestWalkDir(t *testing.T) {
	files, err := pcsutil.WalkDir("/Users/syy/tmp", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
