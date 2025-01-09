package entityconvertors

import (
	"fmt"

	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type binaryItemConvertor struct {
	sic SIConvertor[entities2.StoredItem]
	vc  ValueConvertor[entities2.FileItem, values.FileValue]
}

func (c binaryItemConvertor) ToModel(e *entities2.FileItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("binaryItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c binaryItemConvertor) ToEntity(e *models.StoredItem) (*entities2.FileItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("binaryItemConvertor ToEntity: %w", err)
	}
	return &entities2.FileItem{
		StoredItem: *si,
		FileValue:  *v,
	}, nil
}
