package state

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type AppState struct {
	User       *entities.User
	DBProvider dbProvider
	Storages   *Storages
}
