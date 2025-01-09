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
	User     Storage[entities.User]
	Password Storage[entities.PasswordItem]
	Note     Storage[entities.NoteItem]
	BankCard Storage[entities.BankCardItem]
}

type QueueRunner interface {
	AddJob(queue.Job) error
}
