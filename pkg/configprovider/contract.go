package configprovider

type parser[T any] interface {
	Parse(c *T) error
}
