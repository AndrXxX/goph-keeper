package controllers

import (
	"context"
	"io"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type usersStorage interface {
	Insert(ctx context.Context, m *models.User) (*models.User, error)
	QueryOne(ctx context.Context, login string) (*models.User, error)
}

type hashGenerator interface {
	Generate(data []byte) string
}

type tokenService interface {
	Decrypt(token string) (userID uint, err error)
	Encrypt(userID uint) (token string, err error)
}

type fetcher[T any] interface {
	Fetch(r io.Reader) (*T, error)
}
