package entities

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type NoteItem struct {
	StoredItem
	values.NoteValue
}

func (i NoteItem) FilterValue() string {
	return i.Text + i.Desc
}

func (i NoteItem) Title() string {
	return string(i.Text[0:min(10, len(i.Text))]) + " ..."
}

func (i NoteItem) Description() string {
	return i.Desc
}
