package synchronize

import (
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/synchronizers"
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
	Password synchronizers.Storage[entities.PasswordItem]
	Note     synchronizers.Storage[entities.NoteItem]
	BankCard synchronizers.Storage[entities.BankCardItem]
}

type tokenRefresher func()
