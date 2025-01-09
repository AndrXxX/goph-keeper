package ormmodels

import (
	"time"

	"github.com/google/uuid"
)

type StoredItem struct {
	ID        uuid.UUID
	Desc      string
	UpdatedAt time.Time
}
