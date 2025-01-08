package entities

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type BankCardItem struct {
	StoredItem
	values.BankCardValue
}

func (i BankCardItem) FilterValue() string {
	return i.Number + i.Desc
}

func (i BankCardItem) Title() string {
	return i.Number
}

func (i BankCardItem) Description() string {
	return i.Desc
}
