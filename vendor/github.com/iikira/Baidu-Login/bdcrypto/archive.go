package bdcrypto

import (
	"compress/gzip"
	"io"
	"os"
)

// GZIPCompress GZIP 压缩
func GZIPCompress(src io.Reader, writeTo io.Writer) (err error) {
	w := gzip.NewWriter(writeTo)
	_, err = io.Copy(w, src)
	if err != nil {
		return
	}

	w.Flush()
	return w.Close()
}

// GZIPUncompress GZIP 解压缩
func GZIPUncompress(src io.Reader, writeTo io.Writer) (err error) {
	unReader, err := gzip.NewReader(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(writeTo, unReader)
	if err != nil {
		return
	}

	return unReader.Close()
}

// GZIPCompressFile GZIP 压缩文件
func GZIPCompressFile(filePath string) (err error) {
	return gzipCompressFile("en", filePath)
}

// GZIPUnompressFile GZIP 解压缩文件
func GZIPUnompressFile(filePath string) (err error) {
	return gzipCompressFile("de", filePath)
}

func gzipCompressFile(op, filePath string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return
	}

	tempFilePath := filePath + ".gzip.tmp"
	// 保留文件权限
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fInfo.Mode())
	if err != nil {
		return
	}

	defer tempFile.Close()

	switch op {
	case "en":
		err = GZIPCompress(f, tempFile)
	case "de":
		err = GZIPUncompress(f, tempFile)
	default:
		panic("unknown op" + op)
	}

	if err != nil {
		os.Remove(tempFilePath)
		return
	}

	return os.Rename(tempFilePath, filePath)
}
