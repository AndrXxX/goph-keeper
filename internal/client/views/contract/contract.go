package contract

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
)

type SyncManager interface {
	Sync(dataType string, updates []any) error
}

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}

type Storages struct {
	Password Storage[entities.PasswordItem]
	Note     Storage[entities.NoteItem]
	BankCard Storage[entities.BankCardItem]
}

type QueueRunner interface {
	AddJob(queue.Job) error
}

type UserAccessor interface {
	Auth() error
	SetMasterPass(mp string)
	SetUser(user *entities.User)
}

type BuildInfo struct {
	Version string
	Date    string
}
