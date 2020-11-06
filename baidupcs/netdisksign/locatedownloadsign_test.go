package netdisksign_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs/netdisksign"
	"testing"
)

func TestLocateDownloadSign(t *testing.T) {
	sign := netdisksign.NewLocateDownloadSignWithTimeAndDevUID(1571140066, "O|1E67351CCE80B2CF48DB511CD77ACD9F", 10086, "test_bduss")
	fmt.Printf("%#v\n", sign.Rand == "b6bb7a6f46899e181baea58798d4fdb889775c2c")
}
