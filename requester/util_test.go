package requester

import (
	"strings"
	"testing"
)

// TestParseCookieStrStripQuotes 验证百度部分 cookie 值两侧的双引号 (如 RT) 会被去除,
// 避免 net/http 在发送 Cookie 时报 "invalid byte '"' in Cookie.Value".
func TestParseCookieStrStripQuotes(t *testing.T) {
	cookieStr := `BAIDUID=767874F83F19B61EA4AE1AB875BD9994:FG=1; BDUSS=mhrNmU3ZTFu; STOKEN=e334e904; RT="z=1&dm=baidu.com&si=8da0aa22&hd=1xlf"; csrfToken=5lb_cUD5`

	cookies := ParseCookieStr(cookieStr)

	var rtValue string
	found := false
	for _, c := range cookies {
		if c.Name == "RT" {
			rtValue = c.Value
			found = true
		}
	}

	if !found {
		t.Fatal("RT cookie 未解析出来")
	}
	if strings.Contains(rtValue, "\"") {
		t.Fatalf("RT 值仍包含双引号: %q", rtValue)
	}
	if rtValue != "z=1&dm=baidu.com&si=8da0aa22&hd=1xlf" {
		t.Fatalf("RT 值解析错误: %q", rtValue)
	}

	// 普通无引号 cookie 应原样保留
	for _, c := range cookies {
		if c.Name == "BDUSS" && c.Value != "mhrNmU3ZTFu" {
			t.Fatalf("BDUSS 值错误: %q", c.Value)
		}
	}
}
