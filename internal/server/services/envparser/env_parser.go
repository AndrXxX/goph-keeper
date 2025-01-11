package envparser

import (
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
)

type conf struct {
	Host             string `env:"HOST"`
	DatabaseURI      string `env:"DATABASE_URI"`
	AuthKey          string `env:"AUTH_SECRET_KEY"`
	AuthKeyExpired   int    `env:"AUTH_SECRET_KEY_EXPIRED"`
	PasswordKey      string `env:"PASSWORD_SECRET_KEY"`
	PrivateCryptoKey string `env:"PRIVATE_CRYPTO_KEY"`
	PublicCryptoKey  string `env:"PUBLIC_CRYPTO_KEY"`
}

// Parser сервис для парсинга переменных окружения
type Parser struct {
}

// Parse парсит переменные окружения и наполняет конфигурацию
func (p Parser) Parse(c *config.Config) error {
	cfg := conf{
		Host:             c.Host,
		DatabaseURI:      c.DatabaseURI,
		AuthKey:          c.AuthKey,
		AuthKeyExpired:   c.AuthKeyExpired,
		PasswordKey:      c.PasswordKey,
		PrivateCryptoKey: c.PrivateCryptoKey,
		PublicCryptoKey:  c.PublicCryptoKey,
	}
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("error on parse config: %w", err)
	}
	c.Host = cfg.Host
	c.DatabaseURI = cfg.DatabaseURI
	c.AuthKey = cfg.AuthKey
	c.AuthKeyExpired = cfg.AuthKeyExpired
	c.PasswordKey = cfg.PasswordKey
	c.PrivateCryptoKey = cfg.PrivateCryptoKey
	c.PublicCryptoKey = cfg.PublicCryptoKey
	return nil
}
