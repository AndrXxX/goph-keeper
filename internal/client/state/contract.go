package state

import "gorm.io/gorm"

type dbProvider interface {
	IsDBExist() bool
	DB(masterPass string) (*gorm.DB, error)
	RemoveDB() error
}
