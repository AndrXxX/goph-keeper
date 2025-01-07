package requestsender

import (
	"io"
	"net/http"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type hashGenerator interface {
	Generate(data []byte) string
}

type dataCompressor interface {
	Compress(in []byte) (io.Reader, error)
}

type Option func(*dto.ParamsDto) error
