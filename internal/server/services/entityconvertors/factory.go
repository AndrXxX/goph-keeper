package entityconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type Factory struct {
}

func (f Factory) Password(vc ValueConvertor[entities2.PasswordItem, values.PasswordValue]) Convertor[entities2.PasswordItem] {
	return loginPassItemConvertor{sic: f.StoredItem(datatypes.Passwords), vc: vc}
}

func (f Factory) Note(vc ValueConvertor[entities2.NoteItem, values.NoteValue]) Convertor[entities2.NoteItem] {
	return textItemConvertor{sic: f.StoredItem(datatypes.Notes), vc: vc}
}

func (f Factory) BankCard(vc ValueConvertor[entities2.BankCardItem, values.BankCardValue]) Convertor[entities2.BankCardItem] {
	return bankCardItemConvertor{sic: f.StoredItem(datatypes.BankCards), vc: vc}
}

func (f Factory) File(vc ValueConvertor[entities2.FileItem, values.FileValue]) Convertor[entities2.FileItem] {
	return binaryItemConvertor{sic: f.StoredItem(datatypes.Files), vc: vc}
}

func (f Factory) StoredItem(t string) SIConvertor[entities2.StoredItem] {
	return storedItemConvertor{Type: t}
}
