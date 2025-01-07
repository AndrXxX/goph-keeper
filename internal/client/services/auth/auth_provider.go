package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AndrXxX/goph-keeper/internal/enums/contenttypes"
	"github.com/AndrXxX/goph-keeper/pkg/entities"
)

const (
	loginUrl    = "/api/user/login"
	registerUrl = "/api/user/register"
)

var (
	ExistUser = fmt.Errorf("такой пользователь уже есть")
	WrongCred = fmt.Errorf("неверные логин или пароль")
)

type Provider struct {
	Sender requestSender
	UB     urlBuilder
}

func (p *Provider) Register(u *entities.User) (string, error) {
	return p.send(u, p.UB.Build(registerUrl, nil))
}

func (p *Provider) Login(u *entities.User) (string, error) {
	return p.send(u, p.UB.Build(loginUrl, nil))
}

func (p *Provider) send(u *entities.User, url string) (string, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return "", fmt.Errorf("marshal user: %v", err)
	}
	resp, err := p.Sender.Post(url, contenttypes.ApplicationJSON, data)
	if err != nil {
		return "", fmt.Errorf("send request: %v", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		token := p.getTokenFromHeaders(resp)
		if token == "" {
			return "", fmt.Errorf("token is empty")
		}
		return token, nil
	case http.StatusUnauthorized:
		return "", WrongCred
	case http.StatusConflict:
		return "", ExistUser
	}
	return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
}

func (p *Provider) getTokenFromHeaders(r *http.Response) string {
	if raw := r.Header.Get("Authorization"); raw != "" {
		vals := strings.Split(raw, " ")
		if len(vals) == 2 {
			return vals[1]
		}
	}
	return ""
}
