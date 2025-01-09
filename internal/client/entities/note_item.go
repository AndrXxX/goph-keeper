package entities

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

type NoteItem struct {
	StoredItem
	Text string `json:"text" valid:"required"`
}

func (i NoteItem) FilterValue() string {
	return i.Text + i.Desc
}

func (i NoteItem) Title() string {
	return string(i.Text[0:min(10, len(i.Text))]) + " ..."
}

func (i NoteItem) Description() string {
	if len(i.Desc) > 0 {
		return fmt.Sprintf("%s [%s]", i.Desc, i.UpdatedAt.Format(formats.FullDate))
	}
	return i.UpdatedAt.Format(formats.FullDate)
}
