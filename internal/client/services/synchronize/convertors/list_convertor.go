package convertors

import "fmt"

type ListConvertor[T any] struct{}

func (c ListConvertor[T]) Convert(in []any) ([]T, error) {
	list := make([]T, len(in))
	for i := range in {
		item, ok := in[i].(T)
		if !ok {
			return nil, fmt.Errorf("item has not valid type")
		}
		list[i] = item
	}
	return list, nil
}
