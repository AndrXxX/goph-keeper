package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/pkg/entities"
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type bankCardValueConvertor struct {
}

func (c bankCardValueConvertor) ToValue(v string) (*values.BankCardValue, error) {
	var converted values.BankCardValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal bank card value: %w", err)
	}
	return &converted, nil
}

func (c bankCardValueConvertor) ToString(item *entities.BankCardItem) (string, error) {
	val, err := json.Marshal(values.BankCardValue{
		Number: item.Number, CVCCode: item.CVCCode,
		Validity: item.Validity, Cardholder: item.Cardholder,
	})
	if err != nil {
		return "", fmt.Errorf("marshal bank card value: %w", err)
	}
	return string(val), nil
}
