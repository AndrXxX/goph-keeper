package entities

import (
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type NoteItem struct {
	StoredItem
	values.NoteValue
}
