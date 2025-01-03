package entityconvertors

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type binaryItemConvertor struct {
	sic SIConvertor[entities.StoredItem]
	vc  ValueConvertor[entities.BinaryItem, values.BinaryValue]
}

func (c binaryItemConvertor) ToModel(e *entities.BinaryItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("binaryItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c binaryItemConvertor) ToEntity(e *models.StoredItem) (*entities.BinaryItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("binaryItemConvertor ToEntity: %w", err)
	}
	return &entities.BinaryItem{
		StoredItem:  *si,
		BinaryValue: *v,
	}, nil
}
