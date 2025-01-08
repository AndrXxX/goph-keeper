package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type BankCardOrmEntityConvertor struct {
}

func (c BankCardOrmEntityConvertor) Convert(m *ormmodels.BankCardItem) *entities.BankCardItem {
	return &entities.BankCardItem{
		StoredItem: *ItemOrmEntityConvertor{}.Convert(&m.StoredItem),
		Number:     m.Number,
		CVCCode:    m.CVCCode,
		Validity:   m.Validity,
		Cardholder: m.Cardholder,
	}
}
