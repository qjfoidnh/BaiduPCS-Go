package jsonhelper

import (
	"github.com/json-iterator/go"
	"io"
)

// UnmarshalData 将 r 中的 json 格式的数据, 解析到 data
func UnmarshalData(r io.Reader, data interface{}) error {
	d := jsoniter.NewDecoder(r)
	return d.Decode(data)
}

// MarshalData 将 data, 生成 json 格式的数据, 写入 w 中
func MarshalData(w io.Writer, data interface{}) error {
	e := jsoniter.NewEncoder(w)
	return e.Encode(data)
}
