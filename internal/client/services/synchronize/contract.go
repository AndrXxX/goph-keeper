package synchronize

import (
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type Synchronizer interface {
	Sync(updates []any) error
}

type requestSender interface {
	Get(url string, contentType string) (*http.Response, error)
	Post(url string, contentType string, data []byte) (*http.Response, error)
}

type urlBuilder interface {
	Build(endpoint string, params map[string]string) string
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

type tokenRefresher func()
