package ormstorages

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type storagesFactory struct {
}

func Factory() *storagesFactory {
	return &storagesFactory{}
}

func (f *storagesFactory) User(ctx context.Context, db *gorm.DB) Storage[ormmodels.User] {
	return getStorage[ormmodels.User](ctx, db, new(ormmodels.User))
}

func (f *storagesFactory) Password(ctx context.Context, db *gorm.DB) Storage[ormmodels.PasswordItem] {
	return getStorage[ormmodels.PasswordItem](ctx, db, new(ormmodels.PasswordItem))
}

func (f *storagesFactory) Note(ctx context.Context, db *gorm.DB) Storage[ormmodels.NoteItem] {
	return getStorage[ormmodels.NoteItem](ctx, db, new(ormmodels.NoteItem))
}

func (f *storagesFactory) BankCard(ctx context.Context, db *gorm.DB) Storage[ormmodels.BankCardItem] {
	return getStorage[ormmodels.BankCardItem](ctx, db, new(ormmodels.BankCardItem))
}

func getStorage[T interface{}](ctx context.Context, db *gorm.DB, m *T) Storage[T] {
	s := &ormStorage[T]{db}
	err := s.init(ctx, m)
	if err != nil {
		logger.Log.Error("failed to init storage", zap.Error(err))
	}
	return s
}
