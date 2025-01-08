package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type BankCardEntityOrmConvertor struct {
}

func (c BankCardEntityOrmConvertor) Convert(e *entities.BankCardItem) *ormmodels.BankCardItem {
	return &ormmodels.BankCardItem{
		StoredItem: *ItemEntityOrmConvertor{}.Convert(&e.StoredItem),
		Number:     e.Number,
		CVCCode:    e.CVCCode,
		Validity:   e.Validity,
		Cardholder: e.Cardholder,
	}
}
