package entities

import (
	"time"

	"github.com/google/uuid"
)

type StoredItem struct {
	ID        uuid.UUID `json:"id"`
	Desc      string    `json:"description"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (si StoredItem) GetID() uuid.UUID {
	return si.ID
}

func (si StoredItem) IsStored() bool {
	return si.ID != uuid.Nil
}
