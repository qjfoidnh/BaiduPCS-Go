package baidupcs

import (
	"net/http/cookiejar"
	"strings"
)

type list struct{}

// PublicSuffixList baidupcs PublicSuffixList
var PublicSuffixList cookiejar.PublicSuffixList = list{}

func (list) PublicSuffix(domain string) string {
	if strings.HasSuffix(domain, ".baidu.com") {
		return "com"
	}
	return domain
}

func (list) String() string {
	return "baidupcs"
}
