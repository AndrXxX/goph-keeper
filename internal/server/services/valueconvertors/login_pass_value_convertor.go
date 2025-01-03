package valueconvertors

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type loginPassValueConvertor struct {
}

func (c loginPassValueConvertor) ToValue(v string) (*values.LoginPassValue, error) {
	var converted values.LoginPassValue
	err := json.Unmarshal([]byte(v), &converted)
	if err != nil {
		return nil, fmt.Errorf("unmarshal login_pass value: %w", err)
	}
	return &converted, nil
}

func (c loginPassValueConvertor) ToString(item *entities.LoginPassItem) (string, error) {
	val, err := json.Marshal(values.LoginPassValue{Login: item.Login, Password: item.Password})
	if err != nil {
		return "", fmt.Errorf("marshal login_pass value: %w", err)
	}
	return string(val), nil
}
