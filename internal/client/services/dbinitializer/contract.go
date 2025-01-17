package dbinitializer

import "gorm.io/gorm"

type dbProvider interface {
	DB(key string) (*gorm.DB, error)
	RemoveDB() error
}
