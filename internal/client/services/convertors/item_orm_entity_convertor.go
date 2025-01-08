package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type ItemOrmEntityConvertor struct {
}

func (c ItemOrmEntityConvertor) Convert(m *ormmodels.StoredItem) *entities.StoredItem {
	return &entities.StoredItem{ID: m.ID, Desc: m.Desc, UpdatedAt: m.UpdatedAt}
}
