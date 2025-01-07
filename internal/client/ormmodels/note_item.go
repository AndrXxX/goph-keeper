package ormmodels

import (
	"gorm.io/gorm"
)

type NoteItem struct {
	gorm.Model
	StoredItem
	Text string `json:"text"`
}
