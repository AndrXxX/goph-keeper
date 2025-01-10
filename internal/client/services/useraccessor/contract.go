package useraccessor

import "github.com/AndrXxX/goph-keeper/internal/client/entities"

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}

type authSetup func(u *entities.User)
