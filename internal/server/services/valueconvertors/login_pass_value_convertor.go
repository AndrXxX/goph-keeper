package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/pkg/entities"
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type loginPassValueConvertor struct {
}

func (c loginPassValueConvertor) ToValue(v string) (*values.PasswordValue, error) {
	var converted values.PasswordValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal login_pass value: %w", err)
	}
	return &converted, nil
}

func (c loginPassValueConvertor) ToString(item *entities.PasswordItem) (string, error) {
	val, err := json.Marshal(values.PasswordValue{Login: item.Login, Password: item.Password})
	if err != nil {
		return "", fmt.Errorf("marshal login_pass value: %w", err)
	}
	return string(val), nil
}
