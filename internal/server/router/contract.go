package router

import (
	"context"

	"github.com/google/uuid"
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

type itemsStorage interface {
	Insert(ctx context.Context, m *models.StoredItem) (*models.StoredItem, error)
	Update(ctx context.Context, m *models.StoredItem) (*models.StoredItem, error)
	QueryOneById(ctx context.Context, id uuid.UUID) (*models.StoredItem, error)
	Query(ctx context.Context, m *models.StoredItem) ([]models.StoredItem, error)
}

type Storage struct {
	DB ksql.Provider
	US usersStorage
	IS itemsStorage
}
