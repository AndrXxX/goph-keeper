package gzipcompressor

import "io"

type buffer interface {
	io.Reader
	io.Writer
}
