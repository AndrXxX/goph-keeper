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

func (f Factory) Text() ValueConvertor[entities.TextItem, values.TextValue] {
	return &textValueConvertor{}
}
