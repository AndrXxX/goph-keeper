package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type PasswordEntityOrmConvertor struct {
}

func (c PasswordEntityOrmConvertor) Convert(e *entities.PasswordItem) *ormmodels.PasswordItem {
	return &ormmodels.PasswordItem{
		StoredItem: *ItemEntityOrmConvertor{}.Convert(&e.StoredItem),
		Login:      e.Login,
		Password:   e.Password,
	}
}
