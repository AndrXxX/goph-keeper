package dbprovider

import (
	"fmt"

	"github.com/galeone/igor"
	"github.com/pressly/goose/v3"
)

// DBProvider сервис для предоставления соединения с базой данных
type DBProvider struct {
	DSN string
}

// DB возвращает соединение с базой данных
func (p *DBProvider) DB() (*igor.Database, error) {
	if p.DSN == "" {
		return nil, fmt.Errorf("empty Database DSN")
	}

	db, err := igor.Connect(p.DSN)
	if err != nil {
		return nil, fmt.Errorf("connect db %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("error on goose SetDialect %w", err)
	}

	if err := goose.Up(db.DB(), "internal/server/migrations/postgresql"); err != nil {
		return nil, fmt.Errorf("error on up migrations %w", err)
	}

	return db, nil
}
