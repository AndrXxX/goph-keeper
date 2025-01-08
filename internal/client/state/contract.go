package state

import (
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
)

type dbProvider interface {
	IsDBExist() bool
	DB() (*gorm.DB, error)
	RemoveDB() error
}

type Storages struct {
	User     interfaces.Storage[entities.User]
	Password interfaces.Storage[entities.PasswordItem]
	Note     interfaces.Storage[entities.NoteItem]
	BankCard interfaces.Storage[entities.BankCardItem]
}
