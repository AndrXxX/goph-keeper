package app

import (
	"context"

	"github.com/vingarcia/ksql"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type appConfig struct {
	c *config.Config
}

type usersStorage interface {
	Insert(ctx context.Context, m *models.User) (*models.User, error)
	QueryOne(ctx context.Context, login string) (*models.User, error)
}

type Storage struct {
	DB ksql.Provider
	US usersStorage
}
