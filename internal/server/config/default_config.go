package config

import (
	"github.com/AndrXxX/goph-keeper/internal/enums/defaults"
)

// NewConfig возвращает конфигурацию со значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		Host:            defaults.RunAddress,
		LogLevel:        defaults.LogLevel,
		AuthKey:         defaults.AuthKey,
		AuthKeyExpired:  defaults.AuthKeyExpired,
		PasswordKey:     defaults.PasswordKey,
		FileStoragePath: "./data/server/files",
	}
}
