package escaper

import (
	"strings"
)

type (
	// RuneFunc 判断指定rune
	RuneFunc func(r rune) bool
)

// EscapeByRuneFunc 通过runeFunc转义, runeFunc返回真, 则转义
func EscapeByRuneFunc(s string, runeFunc RuneFunc) string {
	if runeFunc == nil {
		return s
	}

	var (
		builder = &strings.Builder{}
		rs      = []rune(s)
	)

	for k := range rs {
		if !runeFunc(rs[k]) {
			builder.WriteRune(rs[k])
			continue
		}

		if k >= 1 && rs[k-1] == '\\' {
			builder.WriteRune(rs[k])
			continue
		}
		builder.WriteString(`\`)
		builder.WriteRune(rs[k])
	}
	return builder.String()
}

// Escape 转义指定的escapeRunes, 在escapeRunes的前面加上一个反斜杠
func Escape(s string, escapeRunes []rune) string {
	return EscapeByRuneFunc(s, func(r rune) bool {
		for k := range escapeRunes {
			if escapeRunes[k] == r {
				return true
			}
		}
		return false
	})
}

// EscapeStrings 转义字符串数组
func EscapeStrings(ss []string, escapeRunes []rune) {
	for k := range ss {
		ss[k] = Escape(ss[k], escapeRunes)
	}
}

// EscapeStringsByRuneFunc 转义字符串数组, 通过runeFunc
func EscapeStringsByRuneFunc(ss []string, runeFunc RuneFunc) {
	for k := range ss {
		ss[k] = EscapeByRuneFunc(ss[k], runeFunc)
	}
}
