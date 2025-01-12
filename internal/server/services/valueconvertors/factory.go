package valueconvertors

import (
	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type Factory struct {
}

func (f Factory) Password() ValueConvertor[entities2.PasswordItem, values.PasswordValue] {
	return &loginPassValueConvertor{}
}

func (f Factory) Note() ValueConvertor[entities2.NoteItem, values.NoteValue] {
	return &textValueConvertor{}
}

func (f Factory) BankCard() ValueConvertor[entities2.BankCardItem, values.BankCardValue] {
	return &bankCardValueConvertor{}
}

func (f Factory) File() ValueConvertor[entities2.FileItem, values.FileValue] {
	return &fileValueConvertor{}
}
