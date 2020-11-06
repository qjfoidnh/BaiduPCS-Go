package rio

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMultiReaderLen(t *testing.T) {
	rd1, rd2 := strings.NewReader("asdkfljalf"), strings.NewReader("---asva sdf")
	multi := MultiReaderLen(rd1, rd2)
	fmt.Println(multi.Len())
	io.Copy(os.Stdout, multi)
}
