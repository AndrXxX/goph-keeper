package valueconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type Factory struct {
}

func (f Factory) Password() ValueConvertor[entities.PasswordItem, values.PasswordValue] {
	return &loginPassValueConvertor{}
}

func (f Factory) Note() ValueConvertor[entities.NoteItem, values.NoteValue] {
	return &textValueConvertor{}
}

func (f Factory) BankCard() ValueConvertor[entities.BankCardItem, values.BankCardValue] {
	return &bankCardValueConvertor{}
}

func (f Factory) File() ValueConvertor[entities.FileItem, values.FileValue] {
	return &fileValueConvertor{}
}
