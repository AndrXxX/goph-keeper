package requestsender

import (
	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

func WithToken(provider tokenProvider) Option {
	return func(p *dto.ParamsDto) error {
		token := provider.GetToken()
		if token != "" {
			p.Headers["Authorization"] = "Bearer " + token
		}
		return nil
	}
}
