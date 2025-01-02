package models

import (
	"time"
)

type StoredItem struct {
	ID          uint       `ksql:"id"`
	CreatedAt   *time.Time `ksql:"created_at"`
	UpdatedAt   *time.Time `ksql:"updated_at"`
	Type        string     `ksql:"type"`
	Description string     `ksql:"description"`
	Value       string     `ksql:"value"`
	UserID      uint       `ksql:"user_id"`
}

func (u StoredItem) TableName() string {
	return "stored_items"
}
