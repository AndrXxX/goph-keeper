package itemsloader

import (
	"io"
	"net/http"

	"github.com/google/uuid"
)

type requestSender interface {
	Get(url string, contentType string) (*http.Response, error)
	Post(url string, contentType string, data io.Reader) (*http.Response, error)
}

type urlBuilder interface {
	Build(endpoint string, params map[string]string) string
}

type sliceFetcher[T any] interface {
	FetchSlice(r io.Reader) ([]T, error)
}

type fileStorage interface {
	Store(src io.Reader, id uuid.UUID) error
	Get(id uuid.UUID) (file io.ReadCloser, err error)
	IsExist(id uuid.UUID) bool
}
