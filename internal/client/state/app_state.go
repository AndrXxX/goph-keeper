package state

import (
	"context"

	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type AppState struct {
	User            *ormmodels.User
	DBProvider      dbProvider
	Storages        *Storages
	StorageProvider storageProvider
}

func (s *AppState) InitStorages(ctx context.Context, db *gorm.DB) {
	s.Storages = &Storages{
		User:     s.StorageProvider.User(ctx, db),
		Password: s.StorageProvider.Password(ctx, db),
		Note:     s.StorageProvider.Note(ctx, db),
		BankCard: s.StorageProvider.BankCard(ctx, db),
	}
}
