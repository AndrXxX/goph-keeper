package synchronizers

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type loader[T any] interface {
	Download(itemType string) (statusCode int, l []T)
	Upload(itemType string, list []T) (statusCode int, err error)
}

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
}

type fileLoader interface {
	Download() (statusCode int, l []entities.FileItem)
	Upload(list []entities.FileItem) (statusCode int, err error)
}
