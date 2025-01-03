package entityconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type Factory struct {
}

func (f Factory) LoginPass(vc ValueConvertor[entities.LoginPassItem, values.LoginPassValue]) Convertor[entities.LoginPassItem] {
	return loginPassItemConvertor{
		sic: f.StoredItem(),
		vc:  vc,
	}
}

func (f Factory) StoredItem() SIConvertor[entities.StoredItem] {
	return storedItemConvertor{}
}
