package entities

import "github.com/AndrXxX/goph-keeper/internal/server/entities/values"

type TextItem struct {
	StoredItem
	values.TextValue
}
