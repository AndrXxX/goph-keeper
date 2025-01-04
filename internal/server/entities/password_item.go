package entities

import "github.com/AndrXxX/goph-keeper/internal/server/entities/values"

type PasswordItem struct {
	StoredItem
	values.PasswordValue
}
