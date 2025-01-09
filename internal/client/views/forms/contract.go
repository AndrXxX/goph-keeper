package forms

import "github.com/AndrXxX/goph-keeper/internal/client/entities"

type Registerer interface {
	Register(u *entities.User) (string, error)
}

type Loginer interface {
	Login(u *entities.User) (string, error)
}
