package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type PasswordOrmEntityConvertor struct {
}

func (c PasswordOrmEntityConvertor) Convert(m *ormmodels.PasswordItem) *entities.PasswordItem {
	return &entities.PasswordItem{
		StoredItem: *ItemOrmEntityConvertor{}.Convert(&m.StoredItem),
		Login:      m.Login,
		Password:   m.Password,
	}
}
