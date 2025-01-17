package helpers

type Storage[T any] interface {
	FindAll(*T) []T
}
