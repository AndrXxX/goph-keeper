package postgressql

import (
	"context"
	"fmt"

	"github.com/vingarcia/ksql"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

var usersTable = ksql.NewTable("users", "id")

type usersStorage struct {
	db ksql.Provider
}

// Insert вставляет запись
func (s *usersStorage) Insert(ctx context.Context, m *models.User) (*models.User, error) {
	err := s.db.Insert(ctx, usersTable, m)
	if err != nil {
		return nil, fmt.Errorf("insert user %w", err)
	}
	return m, nil
}

// QueryOne извлекает одну запись
func (s *usersStorage) QueryOne(ctx context.Context, login string) (*models.User, error) {
	res := models.User{}
	err := s.db.QueryOne(ctx, &res, "FROM users WHERE login = $1", login)
	if err != nil {
		return nil, fmt.Errorf("queryOne user %w", err)
	}
	return &res, err
}
