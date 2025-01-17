package gzipcompressor

import (
	"compress/gzip"
	"fmt"
	"io"
)

type GzipCompressor struct {
	Buff buffer
}

func (c GzipCompressor) Compress(data io.Reader) (io.Reader, error) {
	if data == nil {
		return data, nil
	}
	w := gzip.NewWriter(c.Buff)
	_, err := io.Copy(w, data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return c.Buff, nil
}
