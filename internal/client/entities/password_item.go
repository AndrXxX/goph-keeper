package entities

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

type PasswordItem struct {
	StoredItem
	Login    string `json:"login" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (i PasswordItem) FilterValue() string {
	return i.Login + i.Desc
}

func (i PasswordItem) Title() string {
	return i.Login
}

func (i PasswordItem) Description() string {
	if len(i.Desc) > 0 {
		return fmt.Sprintf("%s [%s]", i.Desc, i.UpdatedAt.Format(formats.FullDate))
	}
	return i.UpdatedAt.Format(formats.FullDate)
}
