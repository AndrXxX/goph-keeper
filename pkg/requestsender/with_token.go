package requestsender

import (
	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

func WithToken(token string) Option {
	return func(p *dto.ParamsDto) error {
		if token != "" {
			p.Headers["Authorization"] = "Bearer " + token
		}
		return nil
	}
}
