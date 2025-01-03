package valueconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type Factory struct {
}

func (f Factory) LoginPass() ValueConvertor[entities.LoginPassItem, values.LoginPassValue] {
	return &loginPassValueConvertor{}
}
