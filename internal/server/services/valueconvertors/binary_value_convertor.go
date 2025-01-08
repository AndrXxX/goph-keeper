package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type binaryValueConvertor struct {
}

func (c binaryValueConvertor) ToValue(v string) (*values.FileValue, error) {
	var converted values.FileValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal text value: %w", err)
	}
	return &converted, nil
}

func (c binaryValueConvertor) ToString(item *entities.FileItem) (string, error) {
	val, err := json.Marshal(values.FileValue{Data: item.Data})
	if err != nil {
		return "", fmt.Errorf("marshal text value: %w", err)
	}
	return string(val), nil
}
