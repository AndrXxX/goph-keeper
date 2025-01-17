package app

import (
	"context"

	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
)

type queueRunner interface {
	Run(context.Context) error
	Stop(context.Context) error
	AddJob(queue.Job) error
}

type syncManager interface {
	Sync(dataType string, updates []any) error
}

type dbProvider interface {
	DB(key string) (*gorm.DB, error)
	RemoveDB() error
}

type userAccessor interface {
	Auth() error
	SetMasterPass(mp string)
	SetUser(user *entities.User)
	GetUser() *entities.User
	GetToken() string
	SetToken(t string)
	AfterAuth(f func())
}
