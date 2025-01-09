package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type ItemEntityOrmConvertor struct {
}

func (c ItemEntityOrmConvertor) Convert(e *entities.StoredItem) *ormmodels.StoredItem {
	if e == nil {
		return nil
	}
	return &ormmodels.StoredItem{ID: e.ID, Desc: e.Desc, UpdatedAt: e.UpdatedAt}
}
