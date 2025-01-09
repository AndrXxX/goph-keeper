package synchronizers

type loader[T any] interface {
	Download(itemType string) (statusCode int, l []T)
	Upload(itemType string, list []T) (statusCode int, err error)
}

type Storage[T any] interface {
	Find(*T) *T
	Create(*T) (*T, error)
	Update(*T) error
}
