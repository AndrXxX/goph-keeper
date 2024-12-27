package postgressql

import (
	"github.com/galeone/igor"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

// Factory фабрика для хранилищ моделей
type Factory struct {
	DB *igor.Database
}

// UsersStorage хранилище для модели User
func (f *Factory) UsersStorage() *Storage[*models.User] {
	return &Storage[*models.User]{f.DB}
}
