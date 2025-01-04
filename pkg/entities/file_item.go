package entities

import (
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type FileItem struct {
	StoredItem
	values.FileValue
}
