package valueconvertors

type ValueConvertor[ItemT any, ValueT any] interface {
	ToValue(v string) (*ValueT, error)
	ToString(item *ItemT) (string, error)
}
