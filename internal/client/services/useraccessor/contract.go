package useraccessor

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
	FindAll(*T) []T
}

type setupToken func(token string)
type setupDb func(masterPass string) error

type HashGenerator interface {
	Generate(data []byte) string
}

type hashGeneratorFetcher func(key string) HashGenerator
