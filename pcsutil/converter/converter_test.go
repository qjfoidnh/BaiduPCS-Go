package converter_test

import (
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"strings"
	"testing"
)

func TestTrimPathInvalidChars(t *testing.T) {
	trimmed := converter.TrimPathInvalidChars("ksjadfi*/?adf")
	if strings.Compare(trimmed, "ksjadfiadf") != 0 {
		t.Fatalf("trimmed: %s\n", trimmed)
	}
	return
}
