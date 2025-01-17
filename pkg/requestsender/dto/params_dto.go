package dto

import (
	"io"
	"net/http"
)

type ParamsDto struct {
	Buf      io.Reader
	Data     []byte
	Headers  map[string]string
	Response *http.Response
}
