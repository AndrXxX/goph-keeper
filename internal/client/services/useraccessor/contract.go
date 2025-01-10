package useraccessor

import "gorm.io/gorm"

type storageProvider[T any] func(db *gorm.DB) Storage[T]

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}

type setupToken func(token string)
type setupDb func(masterPass string) (*gorm.DB, error)

type HashGenerator interface {
	Generate(data []byte) string
}

type hashGeneratorFetcher func(key string) HashGenerator
