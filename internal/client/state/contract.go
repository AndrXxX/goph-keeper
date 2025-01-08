package state

import (
	"context"

	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type dbProvider interface {
	IsDBExist() bool
	DB() (*gorm.DB, error)
	RemoveDB() error
}

type Storages struct {
	User     interfaces.Storage[ormmodels.User]
	Password interfaces.Storage[ormmodels.PasswordItem]
	Note     interfaces.Storage[ormmodels.NoteItem]
	BankCard interfaces.Storage[ormmodels.BankCardItem]
}

type storageProvider interface {
	User(context.Context, *gorm.DB) interfaces.Storage[ormmodels.User]
	Password(context.Context, *gorm.DB) interfaces.Storage[ormmodels.PasswordItem]
	Note(context.Context, *gorm.DB) interfaces.Storage[ormmodels.NoteItem]
	BankCard(context.Context, *gorm.DB) interfaces.Storage[ormmodels.BankCardItem]
}
