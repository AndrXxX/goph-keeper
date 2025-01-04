package valueconvertors

import (
	entities "github.com/AndrXxX/goph-keeper/pkg/entities"
	values2 "github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type Factory struct {
}

func (f Factory) Password() ValueConvertor[entities.PasswordItem, values2.PasswordValue] {
	return &loginPassValueConvertor{}
}

func (f Factory) Note() ValueConvertor[entities.NoteItem, values2.NoteValue] {
	return &textValueConvertor{}
}

func (f Factory) BankCard() ValueConvertor[entities.BankCardItem, values2.BankCardValue] {
	return &bankCardValueConvertor{}
}

func (f Factory) File() ValueConvertor[entities.FileItem, values2.FileValue] {
	return &binaryValueConvertor{}
}
