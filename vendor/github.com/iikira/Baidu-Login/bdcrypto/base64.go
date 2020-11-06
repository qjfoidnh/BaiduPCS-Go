package bdcrypto

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
)

// Base64Encode base64加密
func Base64Encode(raw []byte) []byte {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	encoder.Write(raw)
	encoder.Close()
	return encoded.Bytes()
}

// Base64Decode base64解密
func Base64Decode(raw []byte) []byte {
	var buf bytes.Buffer
	buf.Write(raw)
	decoder := base64.NewDecoder(base64.StdEncoding, &buf)
	decoded, _ := ioutil.ReadAll(decoder)
	return decoded
}
