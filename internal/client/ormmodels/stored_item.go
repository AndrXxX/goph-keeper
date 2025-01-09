package ormmodels

import (
	"time"
)

type StoredItem struct {
	ID        uint
	Desc      string
	UpdatedAt time.Time
}
