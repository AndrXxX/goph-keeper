package state

import (
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type dbProvider interface {
	IsDBExist() bool
	DB() (*gorm.DB, error)
	RemoveDB() error
}

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}

type Storages struct {
	User Storage[entities.User]
}

type authSetup func(u *entities.User)
