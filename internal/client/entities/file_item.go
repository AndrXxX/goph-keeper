package entities

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/formats"
)

type FileItem struct {
	StoredItem
	Name string
	Path string
}

func (i FileItem) FilterValue() string {
	return i.Name
}

func (i FileItem) Title() string {
	return i.Name
}

func (i FileItem) Description() string {
	if len(i.Desc) > 0 {
		return fmt.Sprintf("%s [%s]", i.Desc, i.UpdatedAt.Format(formats.FullDate))
	}
	return i.UpdatedAt.Format(formats.FullDate)
}
