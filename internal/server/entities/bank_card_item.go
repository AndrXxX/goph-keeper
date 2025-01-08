package entities

import (
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
)

type BankCardItem struct {
	StoredItem
	values.BankCardValue
}
