package ormstorages

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type storagesFactory struct {
	db *gorm.DB
}

func Factory(db *gorm.DB) *storagesFactory {
	return &storagesFactory{db}
}

func (f *storagesFactory) UserStorage(ctx context.Context) *ormStorage[ormmodels.User] {
	return getStorage[ormmodels.User](ctx, f.db, new(ormmodels.User))
}

func (f *storagesFactory) PasswordStorage(ctx context.Context) *ormStorage[ormmodels.PasswordItem] {
	return getStorage[ormmodels.PasswordItem](ctx, f.db, new(ormmodels.PasswordItem))
}

func (f *storagesFactory) NoteStorage(ctx context.Context) *ormStorage[ormmodels.NoteItem] {
	return getStorage[ormmodels.NoteItem](ctx, f.db, new(ormmodels.NoteItem))
}

func (f *storagesFactory) BankCardItemStorage(ctx context.Context) *ormStorage[ormmodels.BankCardItem] {
	return getStorage[ormmodels.BankCardItem](ctx, f.db, new(ormmodels.BankCardItem))
}

func getStorage[T interface{}](ctx context.Context, db *gorm.DB, m *T) *ormStorage[T] {
	s := &ormStorage[T]{db}
	err := s.init(ctx, m)
	if err != nil {
		logger.Log.Error("failed to init storage", zap.Error(err))
	}
	return s
}
