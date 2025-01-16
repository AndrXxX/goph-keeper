package config

import (
	"time"

	"github.com/AndrXxX/goph-keeper/internal/enums/defaults"
)

// NewConfig возвращает конфигурацию со значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		Host:            defaults.RunAddress,
		LogLevel:        defaults.LogLevel,
		LogPath:         defaults.LogPath,
		DBPath:          "./data/client.db",
		ServerKeyPath:   "./data/client/public_server.key",
		FileStoragePath: "./data/client/files/",
		QueueWorkersCnt: 5,
		QueueTimeout:    1 * time.Second,
		ShowMsgTimeout:  2 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		SyncInterval:    10 * time.Second,
	}
}
