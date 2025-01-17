package tokenrefresher

import (
	"fmt"
)

type TokenRefresher struct {
	UserAccessor userAccessor
	Loginer      Loginer
	Storage      Storage
}

func (r *TokenRefresher) Refresh() error {
	token, err := r.Loginer.Login(r.UserAccessor.GetUser())
	if err != nil {
		return fmt.Errorf("refresh token: %w", err)
	}
	r.UserAccessor.SetToken(token)
	uErr := r.Storage.Update(r.UserAccessor.GetUser())
	if uErr != nil {
		return fmt.Errorf("update user after refresh token: %w", uErr)
	}
	return nil
}
