package entities

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type PasswordItem struct {
	StoredItem
	values.PasswordValue
}

func (si PasswordItem) FilterValue() string {
	return si.Login + si.Desc
}

func (si PasswordItem) Title() string {
	return si.Login
}

func (si PasswordItem) Description() string {
	return si.Desc
}
