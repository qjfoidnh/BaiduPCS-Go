package pcsverbose

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsutil/pcstime"
	"io"
	"io/ioutil"
)

//PrintReader 输出Reader
func PrintReader(r io.Reader) {
	b, _ := ioutil.ReadAll(r)
	fmt.Printf("%s\n", b)
}

// PrintArgs 输出字符串数组
func PrintArgs(w io.Writer, args ...string) {
	for k, arg := range args {
		io.WriteString(w, fmt.Sprintf("args[%d] = `%s`, ", k, arg))
	}
	w.Write([]byte{'\n'})
}

func TimePrefix() string {
	return "[" + pcstime.BeijingTimeOption("Refer") + "]"
}
