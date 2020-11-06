package requester

import (
	"fmt"
	"net/http"
	"strings"
)

// ParseCookieStr 解析 Cookie 字符串
func ParseCookieStr(cookieStr string) []*http.Cookie {
	rawCookies := strings.SplitN(cookieStr, ";", -1)
	cookies := make([]*http.Cookie, 0, len(rawCookies))

	for _, rawCookie := range rawCookies {
		s2 := strings.SplitN(rawCookie, "=", 2)
		if len(s2) < 2 {
			fmt.Println(s2)
			continue
		}

		s2[0] = strings.TrimSpace(s2[0])
		s2[1] = strings.TrimSpace(s2[1])

		cookies = append(cookies, &http.Cookie{
			Name:  s2[0],
			Value: s2[1],
		})
	}
	return cookies
}
