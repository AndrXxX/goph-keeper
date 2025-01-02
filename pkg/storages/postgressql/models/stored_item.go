package models

import (
	"time"
)

type StoredItem struct {
	ID          uint      `igor:"primary_key"`
	CreatedAt   time.Time `sql:"default:(now() at time zone 'utc')"`
	UpdatedAt   time.Time `sql:"default:(now() at time zone 'utc')"`
	Type        string
	Description string
	Value       string
	UserID      uint
}

func (u StoredItem) TableName() string {
	return "stored_items"
}
