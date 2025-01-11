package config

import (
	"github.com/AndrXxX/goph-keeper/internal/enums/defaults"
)

// NewConfig возвращает конфигурацию со значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		Host:          defaults.RunAddress,
		LogLevel:      defaults.LogLevel,
		LogPath:       defaults.LogPath,
		DBPath:        "./data/client.db",
		ServerKeyPath: "./data/client/public_server.key",
	}
}
