package entityconvertors

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type Factory struct {
}

func (f Factory) LoginPass(vc ValueConvertor[entities.LoginPassItem, values.LoginPassValue]) Convertor[entities.LoginPassItem] {
	return loginPassItemConvertor{sic: f.StoredItem(), vc: vc}
}

func (f Factory) Text(vc ValueConvertor[entities.TextItem, values.TextValue]) Convertor[entities.TextItem] {
	return textItemConvertor{sic: f.StoredItem(), vc: vc}
}

func (f Factory) BankCard(vc ValueConvertor[entities.BankCardItem, values.BankCardValue]) Convertor[entities.BankCardItem] {
	return bankCardItemConvertor{sic: f.StoredItem(), vc: vc}
}

func (f Factory) Binary(vc ValueConvertor[entities.BinaryItem, values.BinaryValue]) Convertor[entities.BinaryItem] {
	return binaryItemConvertor{sic: f.StoredItem(), vc: vc}
}

func (f Factory) StoredItem() SIConvertor[entities.StoredItem] {
	return storedItemConvertor{}
}
