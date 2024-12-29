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
	fl.Parse()
	return nil
}
