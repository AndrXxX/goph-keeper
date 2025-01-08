package synchronize

import (
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
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

type Storages struct {
	User     interfaces.Storage[entities.User]
	Password interfaces.Storage[entities.PasswordItem]
	Note     interfaces.Storage[entities.NoteItem]
	BankCard interfaces.Storage[entities.BankCardItem]
}

type tokenRefresher func()
