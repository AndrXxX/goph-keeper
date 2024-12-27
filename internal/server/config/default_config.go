package config

import (
	"github.com/AndrXxX/goph-keeper/internal/enums/defaults"
)

// NewConfig возвращает конфигурацию со значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		RunAddress:     defaults.RunAddress,
		LogLevel:       defaults.LogLevel,
		DatabaseURI:    "",
		AuthKey:        defaults.AuthKey,
		AuthKeyExpired: defaults.AuthKeyExpired,
		PasswordKey:    defaults.PasswordKey,
	}
}
