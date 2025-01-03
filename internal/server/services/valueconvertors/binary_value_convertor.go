package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type binaryValueConvertor struct {
}

func (c binaryValueConvertor) ToValue(v string) (*values.BinaryValue, error) {
	var converted values.BinaryValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal text value: %w", err)
	}
	return &converted, nil
}

func (c binaryValueConvertor) ToString(item *entities.BinaryItem) (string, error) {
	val, err := json.Marshal(values.BinaryValue{Data: item.Data})
	if err != nil {
		return "", fmt.Errorf("marshal text value: %w", err)
	}
	return string(val), nil
}
