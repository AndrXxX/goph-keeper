package postgressql

import (
	"github.com/galeone/igor"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type Factory struct {
	DB *igor.Database
}

func (f *Factory) UsersStorage() *Storage[models.User] {
	return &Storage[models.User]{f.DB}
}
