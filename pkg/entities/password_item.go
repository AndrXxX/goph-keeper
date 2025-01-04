package entities

import (
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type PasswordItem struct {
	StoredItem
	values.PasswordValue
}
