package netdisksign_test

import (
	"fmt"
	"github.com/qjfoidnh/Baidu-Login/bdcrypto"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs/netdisksign"
	"testing"
)

func TestSign2(t *testing.T) {
	standard := bdcrypto.Base64Decode([]byte("8RxCbsVeSzn2UjxJAAiV9QQs/WetOj2FJUGwjsMG6SgxFMWlLS/U1Q=="))
	fmt.Println("standard,", standard)
	fmt.Printf("standard s %s\n", standard)

	res := netdisksign.Sign2([]rune("e8c7d729eea7b54551aa594f942decbe"), []rune("37dbe07ade9359c1aa70807e847f768c13360ad2"))
	fmt.Println(res)
	fmt.Printf("%s\n", string(res))
	fmt.Println([]byte(string(res)))
	fmt.Println(bdcrypto.Base64Encode([]byte(string(res))))
}
