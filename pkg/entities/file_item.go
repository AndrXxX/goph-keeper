package entities

import (
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type FileItem struct {
	StoredItem
	values.FileValue
}

func (i FileItem) FilterValue() string {
	return i.Desc
}

func (i FileItem) Title() string {
	return i.Desc
}

func (i FileItem) Description() string {
	return ""

}
