package entities

import "github.com/AndrXxX/goph-keeper/internal/server/entities/values"

type BinaryItem struct {
	StoredItem
	values.BinaryValue
}
