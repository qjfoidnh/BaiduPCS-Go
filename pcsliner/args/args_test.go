package args_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsliner/args"
	"testing"
)

func TestParseArgs(t *testing.T) {
	as := args.Parse(`  one two three "double quotes" 'single quotes'  ""   arg\ with\ spaces "\"quotes\" in 'quotes'" '"quotes" in \'quotes'"  "   `)
	for k := range as {
		fmt.Printf("%d: %s|\n", k, as[k])
	}

	as = args.Parse(` cd  英语_800个有趣句子帮你记忆7000个单词_42页.doc`)
	for k := range as {
		fmt.Printf("%d: %s|\n", k, as[k])
	}
}
