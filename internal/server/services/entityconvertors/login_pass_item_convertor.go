package entityconvertors

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type loginPassItemConvertor struct {
	sic SIConvertor[entities.StoredItem]
	vc  ValueConvertor[entities.LoginPassItem, values.LoginPassValue]
}

func (c loginPassItemConvertor) ToModel(e *entities.LoginPassItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("loginPassItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c loginPassItemConvertor) ToEntity(e *models.StoredItem) (*entities.LoginPassItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("loginPassItemConvertor ToEntity: %w", err)
	}
	return &entities.LoginPassItem{
		StoredItem:     *si,
		LoginPassValue: *v,
	}, nil
}
