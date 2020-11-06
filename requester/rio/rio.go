// Package rio rquester io 工具包
package rio

import (
	"io"
)

type (
	// Lener 返回32-bit长度接口
	Lener interface {
		Len() int
	}

	// Lener64 返回64-bit长度接口
	Lener64 interface {
		Len() int64
	}

	// ReaderLen 实现io.Reader和32-bit长度接口
	ReaderLen interface {
		io.Reader
		Lener
	}

	// ReaderLen64 实现io.Reader和64-bit长度接口
	ReaderLen64 interface {
		io.Reader
		Lener64
	}

	// ReaderAtLen64 实现io.ReaderAt和64-bit长度接口
	ReaderAtLen64 interface {
		io.ReaderAt
		Lener64
	}

	// WriterLen64 实现io.Writer和64-bit长度接口
	WriterLen64 interface {
		io.Writer
		Lener64
	}

	// WriteCloserAt 实现io.WriteCloser和io.WriterAt接口
	WriteCloserAt interface {
		io.WriteCloser
		io.WriterAt
	}

	// WriteCloserLen64At 实现rio.WriteCloserAt和64-bit长度接口
	WriteCloserLen64At interface {
		WriteCloserAt
		Lener64
	}
)
