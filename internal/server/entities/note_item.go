package entities

import "github.com/AndrXxX/goph-keeper/internal/server/entities/values"

type NoteItem struct {
	StoredItem
	values.NoteValue
}
