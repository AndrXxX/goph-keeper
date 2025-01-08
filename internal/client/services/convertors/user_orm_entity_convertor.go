package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type UserOrmEntityConvertor struct {
}

func (c UserOrmEntityConvertor) Convert(m *ormmodels.User) *entities.User {
	if m == nil {
		return nil
	}
	return &entities.User{
		ID:             m.ID,
		Login:          m.Login,
		Password:       m.Password,
		Token:          m.Token,
		MasterPassword: m.MasterPassword,
	}
}
