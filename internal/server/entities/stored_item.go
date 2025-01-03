package entities

import "time"

type StoredItem struct {
	ID          uint       `json:"id"`
	Description string     `json:"description"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
