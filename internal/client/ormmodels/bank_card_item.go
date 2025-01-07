package ormmodels

import (
	"gorm.io/gorm"
)

type BankCardItem struct {
	gorm.Model
	StoredItem
	Number     string
	CVCCode    string
	Validity   string
	Cardholder string
}
