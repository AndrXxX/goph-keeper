package entities

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

type PasswordItem struct {
	StoredItem
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (i PasswordItem) FilterValue() string {
	return i.Login + i.Desc
}

func (i PasswordItem) Title() string {
	return i.Login
}

func (i PasswordItem) Description() string {
	return fmt.Sprintf("%s [%s]", i.Desc, i.UpdatedAt.Format(formats.FullDate))
}
