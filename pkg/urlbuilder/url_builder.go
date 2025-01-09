package urlbuilder

import (
	"net/url"
	"strings"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type urlBuilder struct {
	host string
}

func New(host string) *urlBuilder {
	u, err := url.Parse(host)
	if err != nil {
		logger.Log.Error("Error on parse host", zap.String("host", host), zap.Error(err))
		return nil
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Scheme == "localhost" {
		u.Scheme = "http://localhost"
	}
	return &urlBuilder{host: u.String()}
}

func (b *urlBuilder) Build(endpoint string, params map[string]string) string {
	for k, v := range params {
		endpoint = strings.Replace(endpoint, "{"+k+"}", v, -1)
	}
	return b.host + endpoint
}
