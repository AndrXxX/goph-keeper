package requestsender

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

// RequestSender сервис для отправки запросов
type RequestSender struct {
	c    client
	opts []Option
}

// New возвращает сервис RequestSender для отправки запросов
func New(c client, opts ...Option) *RequestSender {
	return &RequestSender{c, opts}
}

// Post отправляет запрос методом Post
func (s *RequestSender) Post(url string, contentType string, data io.Reader) (*http.Response, error) {
	params := dto.ParamsDto{
		Buf:     data,
		Headers: map[string]string{"Content-Type": contentType},
	}
	for _, opt := range s.opts {
		err := opt(&params)
		if err != nil {
			return nil, fmt.Errorf("request sender set options before request: %w", err)
		}
	}
	r, err := http.NewRequest("POST", url, params.Buf)
	if err != nil {
		return nil, fmt.Errorf("request sender creating request: %w", err)
	}
	for k, v := range params.Headers {
		r.Header.Set(k, v)
	}

	resp, err := s.c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("request sender do request: %w", err)
	}
	params.Response = resp
	params.Buf = nil
	for _, opt := range s.opts {
		err := opt(&params)
		if err != nil {
			return nil, fmt.Errorf("request sender set options after request: %w", err)
		}
	}
	return resp, nil
}

// Get отправляет запрос методом Get
func (s *RequestSender) Get(url string, contentType string) (*http.Response, error) {
	params := dto.ParamsDto{
		Headers: map[string]string{"Content-Type": contentType},
	}
	for _, opt := range s.opts {
		err := opt(&params)
		if err != nil {
			return nil, err
		}
	}
	r, err := http.NewRequest("GET", url, params.Buf)
	if err != nil {
		return nil, err
	}
	for k, v := range params.Headers {
		r.Header.Set(k, v)
	}

	resp, err := s.c.Do(r)
	params.Response = resp
	for _, opt := range s.opts {
		err := opt(&params)
		if err != nil {
			return nil, err
		}
	}
	return resp, err
}
