package args

import (
	"strings"
	"unicode"
)

const (
	CharEscape      = '\\'
	CharSingleQuote = '\''
	CharDoubleQuote = '"'
	CharBackQuote   = '`'
)

// IsQuote 是否为引号
func IsQuote(r rune) bool {
	return r == CharSingleQuote || r == CharDoubleQuote || r == CharBackQuote
}

// Parse 解析line, 忽略括号
func Parse(line string) (lineArgs []string) {
	var (
		rl        = []rune(line + " ")
		buf       = strings.Builder{}
		quoteChar rune
		nextChar  rune
		escaped   bool
		in        bool
	)

	var (
		isSpace bool
	)

	for k, r := range rl {
		isSpace = unicode.IsSpace(r)
		if !isSpace && !in {
			in = true
		}

		switch {
		case escaped: // 已转义, 跳过
			escaped = false
			//pass
		case r == CharEscape: // 转义模式
			if k+1+1 < len(rl) { // 不是最后一个字符, 多+1是因为最后一个空格
				nextChar = rl[k+1]
				// 仅支持转义这些字符, 否则原样输出反斜杠
				if unicode.IsSpace(nextChar) || IsQuote(nextChar) || nextChar == CharEscape {
					escaped = true
					continue
				}
			}
			// pass
		case IsQuote(r):
			if quoteChar == 0 { //未引
				quoteChar = r
				continue
			}

			if quoteChar == r { //取消引
				quoteChar = 0
				continue
			}
		case isSpace:
			if !in { // 忽略多余的空格
				continue
			}
			if quoteChar == 0 { // 未在引号内
				lineArgs = append(lineArgs, buf.String())
				buf.Reset()
				in = false
				continue
			}
		}

		buf.WriteRune(r)
	}

	return
}
