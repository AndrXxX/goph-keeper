package entityconvertors

import (
	"fmt"

	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type bankCardItemConvertor struct {
	sic SIConvertor[entities2.StoredItem]
	vc  ValueConvertor[entities2.BankCardItem, values.BankCardValue]
}

func (c bankCardItemConvertor) ToModel(e *entities2.BankCardItem, userID uint) (*models.StoredItem, error) {
	v, err := c.vc.ToString(e)
	if err != nil {
		return nil, fmt.Errorf("bankCardItemConvertor ToModel: %w", err)
	}
	return c.sic.ToModel(&e.StoredItem, userID, v), nil
}

func (c bankCardItemConvertor) ToEntity(e *models.StoredItem) (*entities2.BankCardItem, error) {
	si := c.sic.ToEntity(e)
	v, err := c.vc.ToValue(e.Value)
	if err != nil {
		return nil, fmt.Errorf("bankCardItemConvertor ToEntity: %w", err)
	}
	return &entities2.BankCardItem{
		StoredItem:    *si,
		BankCardValue: *v,
	}, nil
}
