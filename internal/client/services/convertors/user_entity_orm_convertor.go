package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type UserEntityOrmConvertor struct {
}

func (c UserEntityOrmConvertor) Convert(e *entities.User) *ormmodels.User {
	if e == nil {
		return nil
	}
	return &ormmodels.User{
		ID:             e.ID,
		Login:          e.Login,
		Password:       e.Password,
		Token:          e.Token,
		MasterPassword: e.MasterPassword,
	}
}
