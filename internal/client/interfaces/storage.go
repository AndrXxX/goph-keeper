package interfaces

type Storage[T any] interface {
	Find(m *T) *T
	Create(m *T) (*T, error)
	Update(m *T) error
	FindAll(m *T) []T
}
