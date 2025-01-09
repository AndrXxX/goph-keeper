package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	jwt.RegisteredClaims
	UserID uint
}

type tokenService struct {
	key     string
	expired time.Duration
}

// New создает экземпляр сервиса для шифровки/дешифровки токена
func New(key string, expired time.Duration) *tokenService {
	return &tokenService{key, expired}
}

// Encrypt генерирует токен по userID
func (ts *tokenService) Encrypt(userID uint) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.expired)),
		},
		UserID: userID,
	})
	token, err = t.SignedString([]byte(ts.key))
	return token, err
}

// Decrypt извлекает userID из токена
func (ts *tokenService) Decrypt(token string) (userID uint, err error) {
	claims := &claims{}
	t, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(ts.key), nil
		})
	if err != nil {
		return 0, err
	}
	if !t.Valid {
		return 0, fmt.Errorf("token is not valid")
	}
	return claims.UserID, nil
}
