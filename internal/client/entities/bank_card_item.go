package entities

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

type BankCardItem struct {
	StoredItem
	Number     string `json:"number" valid:"required,luhn~Card number is not valid"`
	CVCCode    string `json:"cvc_code" valid:"required,numeric,stringlength(3|3)~CVCCode must contain 3 digits"`
	Validity   string `json:"validity" valid:"required,cardvalidity~Validity must contain date in format 10/2024"`
	Cardholder string `json:"cardholder" valid:"required,minstringlength(10)~Cardholder must contain 10 characters"`
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
