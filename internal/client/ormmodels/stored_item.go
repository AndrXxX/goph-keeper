package ormmodels

import (
	"time"

	"gorm.io/gorm"
)

type StoredItem struct {
	gorm.Model
	ID        uint
	Desc      string
	UpdatedAt *time.Time
}
