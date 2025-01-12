package flagsparser

import (
	fl "flag"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
)

// Parser сервис для парсинга аргументов командной строки
type Parser struct {
}

// Parse парсит аргументы командной строки и наполняет конфигурацию
func (p Parser) Parse(c *config.Config) error {
	fl.StringVar(&c.Host, "a", c.Host, "Net address host:port")
	fl.StringVar(&c.DatabaseURI, "d", c.DatabaseURI, "Database URI")
	fl.StringVar(&c.AuthKey, "ak", c.AuthKey, "Auth Key")
	fl.IntVar(&c.AuthKeyExpired, "ake", c.AuthKeyExpired, "Auth Key expired")
	fl.StringVar(&c.PasswordKey, "pk", c.PasswordKey, "Password Key")
	fl.StringVar(&c.PrivateKeyPath, "privateCK", c.PrivateKeyPath, "Private crypto key path")
	fl.StringVar(&c.PublicKeyPath, "publicCK", c.PublicKeyPath, "Public crypto key path")
	fl.StringVar(&c.FileStoragePath, "fsp", c.FileStoragePath, "File storage path")
	fl.Parse()
	return nil
}
