package views

import "github.com/AndrXxX/goph-keeper/pkg/entities"

type registerer interface {
	Register(u *entities.User) (string, error)
}

type loginer interface {
	Login(u *entities.User) (string, error)
}
