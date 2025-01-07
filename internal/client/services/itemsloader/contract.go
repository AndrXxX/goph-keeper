package itemsloader

import (
	"io"
	"net/http"
)

type requestSender interface {
	Get(url string, contentType string) (*http.Response, error)
	Post(url string, contentType string, data []byte) (*http.Response, error)
}

type urlBuilder interface {
	Build(endpoint string, params map[string]string) string
}

type sliceFetcher[T any] interface {
	FetchSlice(r io.Reader) ([]T, error)
}
