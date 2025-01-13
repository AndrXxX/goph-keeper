package synchronize

import (
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/synchronizers"
)

type Synchronizer interface {
	Sync(updates []any) error
}

type requestSender interface {
	Get(url string, contentType string) (*http.Response, error)
	Post(url string, contentType string, data io.Reader) (*http.Response, error)
}

type urlBuilder interface {
	Build(endpoint string, params map[string]string) string
}

type fileStorage interface {
	Store(src io.Reader, id uuid.UUID) error
	Get(id uuid.UUID) (file io.ReadCloser, err error)
	IsExist(id uuid.UUID) bool
}

type Storages struct {
	Password synchronizers.Storage[entities.PasswordItem]
	Note     synchronizers.Storage[entities.NoteItem]
	BankCard synchronizers.Storage[entities.BankCardItem]
	File     synchronizers.Storage[entities.FileItem]
	FS       fileStorage
}

type tokenRefresher func()
