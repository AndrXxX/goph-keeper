package entities

import (
	"github.com/AndrXxX/goph-keeper/pkg/entities/values"
)

type BankCardItem struct {
	StoredItem
	values.BankCardValue
}
