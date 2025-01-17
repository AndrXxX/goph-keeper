package requestsender

import (
	"bytes"
	"fmt"
	"io"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

func WithSHA256(hg hashGenerator) Option {
	return func(p *dto.ParamsDto) error {
		if p.Buf == nil {
			return nil
		}
		encoded, err := io.ReadAll(p.Buf)
		if err != nil {
			return fmt.Errorf("error on read encoded data: %w", err)
		}
		p.Buf = bytes.NewBuffer(encoded)
		p.Headers["HashSHA256"] = hg.Generate(encoded)
		return nil
	}
}
