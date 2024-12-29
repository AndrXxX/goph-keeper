package controllers

import (
	"io"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type userService interface {
	Find(u *models.User) (*models.User, error)
	Create(u *models.User) (*models.User, error)
}

type hashGenerator interface {
	Generate(data []byte) string
}

type tokenService interface {
	Decrypt(token string) (userID uint, err error)
	Encrypt(userID uint) (token string, err error)
}

type userJSONRequestFetcher interface {
	Fetch(r io.Reader) (*entities.User, error)
}
