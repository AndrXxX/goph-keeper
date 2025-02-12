package entityconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type storedItemConvertor struct {
	Type string
}

func (c storedItemConvertor) ToModel(e *entities.StoredItem, userID uint, v string) *models.StoredItem {
	return &models.StoredItem{
		ID:          e.ID,
		Type:        c.Type,
		UpdatedAt:   &e.UpdatedAt,
		Description: e.Desc,
		Value:       v,
		UserID:      userID,
	}
}

func (c storedItemConvertor) ToEntity(e *models.StoredItem) *entities.StoredItem {
	return &entities.StoredItem{
		ID:        e.ID,
		Desc:      e.Description,
		UpdatedAt: *e.UpdatedAt,
	}
}
