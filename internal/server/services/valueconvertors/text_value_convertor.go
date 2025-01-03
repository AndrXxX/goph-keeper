package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type textValueConvertor struct {
}

func (c textValueConvertor) ToValue(v string) (*values.TextValue, error) {
	var converted values.TextValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal text value: %w", err)
	}
	return &converted, nil
}

func (c textValueConvertor) ToString(item *entities.TextItem) (string, error) {
	val, err := json.Marshal(values.TextValue{Text: item.Text})
	if err != nil {
		return "", fmt.Errorf("marshal text value: %w", err)
	}
	return string(val), nil
}
