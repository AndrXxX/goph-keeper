package entities

import "time"

type StoredItem struct {
	ID        uint       `json:"id"`
	Desc      string     `json:"description"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (si StoredItem) GetID() uint {
	return si.ID
}
