package entityconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	entities "github.com/AndrXxX/goph-keeper/pkg/entities"
	values2 "github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type Factory struct {
}

func (f Factory) Password(vc ValueConvertor[entities.PasswordItem, values2.PasswordValue]) Convertor[entities.PasswordItem] {
	return loginPassItemConvertor{sic: f.StoredItem(datatypes.Passwords), vc: vc}
}

func (f Factory) Note(vc ValueConvertor[entities.NoteItem, values2.NoteValue]) Convertor[entities.NoteItem] {
	return textItemConvertor{sic: f.StoredItem(datatypes.Notes), vc: vc}
}

func (f Factory) BankCard(vc ValueConvertor[entities.BankCardItem, values2.BankCardValue]) Convertor[entities.BankCardItem] {
	return bankCardItemConvertor{sic: f.StoredItem(datatypes.BankCards), vc: vc}
}

func (f Factory) File(vc ValueConvertor[entities.FileItem, values2.FileValue]) Convertor[entities.FileItem] {
	return binaryItemConvertor{sic: f.StoredItem(datatypes.Files), vc: vc}
}

func (f Factory) StoredItem(t string) SIConvertor[entities.StoredItem] {
	return storedItemConvertor{Type: t}
}
