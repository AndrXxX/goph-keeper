package storageadapters

type convertor[T1 any, T2 any] interface {
	Convert(e *T1) *T2
}

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}
