package entityconvertors

import (
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type Convertor[E any] interface {
	ToModel(e *E, userID uint) (*models.StoredItem, error)
	ToEntity(e *models.StoredItem) (*E, error)
}

type SIConvertor[E any] interface {
	ToModel(e *E, userID uint, v string) *models.StoredItem
	ToEntity(e *models.StoredItem) *E
}

type ValueConvertor[ItemT any, ValueT any] interface {
	ToValue(v string) (*ValueT, error)
	ToString(item *ItemT) (string, error)
}
