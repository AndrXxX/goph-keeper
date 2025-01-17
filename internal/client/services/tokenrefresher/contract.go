package tokenrefresher

import "github.com/AndrXxX/goph-keeper/internal/client/entities"

type userAccessor interface {
	GetUser() *entities.User
	SetToken(t string)
}

type Loginer interface {
	Login(u *entities.User) (string, error)
}

type Storage interface {
	Update(user *entities.User) error
}
