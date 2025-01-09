package entityconvertors

import (
	"fmt"

	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type textItemConvertor struct {
	sic SIConvertor[entities2.StoredItem]
	vc  ValueConvertor[entities2.NoteItem, values.NoteValue]
}

func (c textItemConvertor) ToModel(e *entities2.NoteItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("textItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c textItemConvertor) ToEntity(e *models.StoredItem) (*entities2.NoteItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("textItemConvertor ToEntity: %w", err)
	}
	return &entities2.NoteItem{
		StoredItem: *si,
		NoteValue:  *v,
	}, nil
}
