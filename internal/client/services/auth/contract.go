package auth

import (
	"io"
	"net/http"
)

type requestSender interface {
	Get(url string, contentType string) (*http.Response, error)
	Post(url string, contentType string, data io.Reader) (*http.Response, error)
}

type urlBuilder interface {
	Build(endpoint string, params map[string]string) string
}

type keySaver interface {
	Store(resp *http.Response) error
}
