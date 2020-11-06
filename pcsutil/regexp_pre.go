package pcsutil

import (
	"regexp"
)

var (
	// HTTPSRE https regexp
	HTTPSRE = regexp.MustCompile("^https")
	// ChinaPhoneRE https regexp
	ChinaPhoneRE = regexp.MustCompile(`^(\+86)?1[3-9][0-9]\d{8}$`)
)
