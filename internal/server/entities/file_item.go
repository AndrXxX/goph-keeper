package entities

import "github.com/AndrXxX/goph-keeper/internal/server/entities/values"

type FileItem struct {
	StoredItem
	values.FileValue
}
