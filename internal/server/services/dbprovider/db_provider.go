package dbprovider

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	return &db, nil
}
