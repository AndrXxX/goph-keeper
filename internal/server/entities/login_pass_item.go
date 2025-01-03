package entities

import "github.com/AndrXxX/goph-keeper/internal/server/entities/values"

type LoginPassItem struct {
	StoredItem
	values.LoginPassValue
}
