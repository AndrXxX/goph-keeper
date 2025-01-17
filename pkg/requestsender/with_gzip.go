package requestsender

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/pkg/gzipcompressor"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

func WithGzip(comp dataCompressor) Option {
	return func(p *dto.ParamsDto) error {
		if p.Buf == nil {
			return nil
		}
		buf, err := comp.Compress(p.Buf)
		if err != nil {
			return err
		}
		p.Buf = buf
		p.Headers["Content-Encoding"] = "gzip"
		p.Headers["Accept-Encoding"] = "gzip"

		if p.Response == nil {
			return nil
		}
		cr, err := gzipcompressor.NewCompressReader(p.Response.Body)
		if err != nil {
			return fmt.Errorf("creating gzip compressor: %w", err)
		}
		p.Response.Body = cr
		defer func() {
			_ = cr.Close()
		}()
		return nil
	}
}
