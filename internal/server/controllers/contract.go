package controllers

import "github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"

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
