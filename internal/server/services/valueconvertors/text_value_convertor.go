package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/pkg/entities"
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type textValueConvertor struct {
}

func (c textValueConvertor) ToValue(v string) (*values.NoteValue, error) {
	var converted values.NoteValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal text value: %w", err)
	}
	return &converted, nil
}

func (c textValueConvertor) ToString(item *entities.NoteItem) (string, error) {
	val, err := json.Marshal(values.NoteValue{Text: item.Text})
	if err != nil {
		return "", fmt.Errorf("marshal text value: %w", err)
	}
	return string(val), nil
}
