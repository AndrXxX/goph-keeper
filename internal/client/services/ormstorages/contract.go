package ormstorages

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}
