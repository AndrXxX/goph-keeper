package ormmodels

import "gorm.io/gorm"

type PasswordItem struct {
	gorm.Model
	StoredItem
	Login    string
	Password string
}
