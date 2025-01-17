package app

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type tokenRefresher struct {
	ua userAccessor
	l  Loginer
	us Storage[entities.User]
}

func (r *tokenRefresher) Refresh() error {
	token, err := r.l.Login(r.ua.GetUser())
	if err != nil {
		return fmt.Errorf("refresh token: %w", err)
	}
	r.ua.SetToken(token)
	uErr := r.us.Update(r.ua.GetUser())
	if uErr != nil {
		return fmt.Errorf("update user after refresh token: %w", uErr)
	}
	return nil
}
