// Package multipartreader helps you encode large files in MIME multipart format
// without reading the entire content into memory.
package multipartreader

import (
	"errors"
	"fmt"
	"github.com/qjfoidnh/BaiduPCS-Go/requester/rio"
	"io"
	"mime/multipart"
	"strings"
	"sync"
	"sync/atomic"
)

type (
	// MultipartReader MIME multipart format
	MultipartReader struct {
		length      int64
		contentType string
		boundary    string

		formBody  string
		parts     []*part
		part64s   []*part64
		formClose string

		mu          sync.Mutex
		closed      bool
		multiReader io.Reader
	}

	part struct {
		form      string
		readerlen rio.ReaderLen
	}

	part64 struct {
		form        string
		readerlen64 rio.ReaderLen64
	}
)

// NewMultipartReader 返回初始化的 *MultipartReader
func NewMultipartReader() (mr *MultipartReader) {
	builder := &strings.Builder{}
	writer := multipart.NewWriter(builder)
	mr = &MultipartReader{
		contentType: writer.FormDataContentType(),
		boundary:    writer.Boundary(),
	}

	mr.length += int64(builder.Len())
	mr.formBody = builder.String()
	return
}

// AddFormField 增加 form 表单
func (mr *MultipartReader) AddFormField(fieldname string, readerlen rio.ReaderLen) {
	if readerlen == nil {
		return
	}

	mpart := &part{
		form:      fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"%s\"\r\n\r\n", mr.boundary, fieldname),
		readerlen: readerlen,
	}
	atomic.AddInt64(&mr.length, int64(len(mpart.form)+mpart.readerlen.Len()))
	mr.parts = append(mr.parts, mpart)
}

// AddFormFile 增加 form 文件表单
func (mr *MultipartReader) AddFormFile(fieldname, filename string, readerlen64 rio.ReaderLen64) {
	if readerlen64 == nil {
		return
	}

	mpart64 := &part64{
		form:        fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"%s\"; filename=\"%s\"\r\n\r\n", mr.boundary, fieldname, filename),
		readerlen64: readerlen64,
	}
	atomic.AddInt64(&mr.length, int64(len(mpart64.form))+mpart64.readerlen64.Len())
	mr.part64s = append(mr.part64s, mpart64)
}

//CloseMultipart 关闭multipartreader
func (mr *MultipartReader) CloseMultipart() error {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	if mr.closed {
		return errors.New("multipartreader already closed")
	}

	mr.formClose = "\r\n--" + mr.boundary + "--\r\n"
	atomic.AddInt64(&mr.length, int64(len(mr.formClose)))

	numReaders := 0
	if mr.formBody != "" {
		numReaders++
	}
	numReaders += 2*len(mr.parts) + 2*len(mr.part64s)
	if mr.formClose != "" {
		numReaders++
	}

	readers := make([]io.Reader, 0, numReaders)
	readers = append(readers, strings.NewReader(mr.formBody))
	for k := range mr.parts {
		readers = append(readers, strings.NewReader(mr.parts[k].form), mr.parts[k].readerlen)
	}
	for k := range mr.part64s {
		readers = append(readers, strings.NewReader(mr.part64s[k].form), mr.part64s[k].readerlen64)
	}
	readers = append(readers, strings.NewReader(mr.formClose))
	mr.multiReader = io.MultiReader(readers...)

	mr.closed = true
	return nil
}

//ContentType 返回Content-Type
func (mr *MultipartReader) ContentType() string {
	return mr.contentType
}

func (mr *MultipartReader) Read(p []byte) (n int, err error) {
	if !mr.closed {
		return 0, errors.New("multipartreader not closed")
	}
	n, err = mr.multiReader.Read(p)
	return n, err
}

// Len 返回表单内容总长度
func (mr *MultipartReader) Len() int64 {
	return atomic.LoadInt64(&mr.length)
}
