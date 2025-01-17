package requestsender

import (
	"io"
	"net/http"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

type tokenProvider interface {
	GetToken() string
}

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type hashGenerator interface {
	Generate(data []byte) string
}

type dataCompressor interface {
	Compress(in io.Reader) (io.Reader, error)
}

type Option func(*dto.ParamsDto) error
