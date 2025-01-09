package entities

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

type BankCardItem struct {
	StoredItem
	Number     string `json:"number"`
	CVCCode    string `json:"cvc_code"`
	Validity   string `json:"validity"`
	Cardholder string `json:"cardholder"`
}

func (i BankCardItem) FilterValue() string {
	return i.Number + i.Desc
}

func (i BankCardItem) Title() string {
	return i.Number
}

func (i BankCardItem) Description() string {
	if len(i.Desc) > 0 {
		return fmt.Sprintf("%s [%s]", i.Desc, i.UpdatedAt.Format(formats.FullDate))
	}
	return i.UpdatedAt.Format(formats.FullDate)
}
