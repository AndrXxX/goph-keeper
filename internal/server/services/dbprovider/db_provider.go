package dbprovider

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/vingarcia/ksql"
	kpgx "github.com/vingarcia/ksql/adapters/kpgx5"
)

// DBProvider сервис для предоставления соединения с базой данных
type DBProvider struct {
	DSN string
}

// DB возвращает соединение с базой данных
func (p *DBProvider) DB(ctx context.Context) (ksql.Provider, error) {
	if p.DSN == "" {
		return nil, fmt.Errorf("empty Database DSN")
	}

	db, err := kpgx.New(ctx, p.DSN, ksql.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect db %w", err)
	}
	return &db, p.migrate()
}

func (p *DBProvider) migrate() error {
	db, err := sql.Open("pgx", p.DSN)
	if err != nil {
		return fmt.Errorf("error opening db %w", err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error on goose SetDialect %w", err)
	}
	if err := goose.Up(db, "internal/server/migrations/postgresql"); err != nil {
		return fmt.Errorf("error on up migrations %w", err)
	}
	return nil
}
