package entityconvertors

import (
	"fmt"

	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type loginPassItemConvertor struct {
	sic SIConvertor[entities2.StoredItem]
	vc  ValueConvertor[entities2.PasswordItem, values.PasswordValue]
}

func (c loginPassItemConvertor) ToModel(e *entities2.PasswordItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("loginPassItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c loginPassItemConvertor) ToEntity(e *models.StoredItem) (*entities2.PasswordItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("loginPassItemConvertor ToEntity: %w", err)
	}
	return &entities2.PasswordItem{
		StoredItem:    *si,
		PasswordValue: *v,
	}, nil
}
