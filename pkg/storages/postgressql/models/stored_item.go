package models

import (
	"time"

	"github.com/google/uuid"
)

type StoredItem struct {
	ID          uuid.UUID  `ksql:"id"`
	CreatedAt   *time.Time `ksql:"created_at"`
	UpdatedAt   *time.Time `ksql:"updated_at"`
	Type        string     `ksql:"type"`
	Description string     `ksql:"description"`
	Value       string     `ksql:"value"`
	UserID      uint       `ksql:"user_id"`
}
