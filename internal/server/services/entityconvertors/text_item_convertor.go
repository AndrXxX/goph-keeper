package entityconvertors

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type textItemConvertor struct {
	sic SIConvertor[entities.StoredItem]
	vc  ValueConvertor[entities.TextItem, values.TextValue]
}

func (c textItemConvertor) ToModel(e *entities.TextItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("textItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c textItemConvertor) ToEntity(e *models.StoredItem) (*entities.TextItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("textItemConvertor ToEntity: %w", err)
	}
	return &entities.TextItem{
		StoredItem: *si,
		TextValue:  *v,
	}, nil
}
