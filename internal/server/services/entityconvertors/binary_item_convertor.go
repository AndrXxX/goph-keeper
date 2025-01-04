package entityconvertors

import (
	"fmt"

	entities "github.com/AndrXxX/goph-keeper/pkg/entities"
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type binaryItemConvertor struct {
	sic SIConvertor[entities.StoredItem]
	vc  ValueConvertor[entities.FileItem, values.FileValue]
}

func (c binaryItemConvertor) ToModel(e *entities.FileItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("binaryItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c binaryItemConvertor) ToEntity(e *models.StoredItem) (*entities.FileItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("binaryItemConvertor ToEntity: %w", err)
	}
	return &entities.FileItem{
		StoredItem: *si,
		FileValue:  *v,
	}, nil
}
