package entities

import "github.com/AndrXxX/goph-keeper/internal/client/formats"

type FileItem struct {
	StoredItem
	Data string `json:"data"`
}

func (i FileItem) FilterValue() string {
	return i.Desc
}

func (i FileItem) Title() string {
	return i.Desc
}

func (i FileItem) Description() string {
	return i.UpdatedAt.Format(formats.FullDate)
}
