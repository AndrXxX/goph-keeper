package app

import (
	"github.com/galeone/igor"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type appConfig struct {
	c *config.Config
}

type userService interface {
	Find(u *models.User) (*models.User, error)
	Create(u *models.User) (*models.User, error)
}

type Storage struct {
	DB *igor.Database
	US userService
}
