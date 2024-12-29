package configprovider

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type configProvider[T any] struct {
	parsers       []parser[T]
	defaultConfig *T
}

// Fetch извлекает, проверяет и возвращает конфигурацию
func (p *configProvider[T]) Fetch() (*T, error) {
	for _, pr := range p.parsers {
		if err := pr.Parse(p.defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
	}
	if _, err := govalidator.ValidateStruct(p.defaultConfig); err != nil {
		return nil, fmt.Errorf("failed to validate env vars: %w", err)
	}
	return p.defaultConfig, nil
}

// New возвращает сервис configProvider для извлечения конфигурации
func New[T any](def *T, parsers ...parser[T]) *configProvider[T] {
	return &configProvider[T]{parsers: parsers, defaultConfig: def}
}
