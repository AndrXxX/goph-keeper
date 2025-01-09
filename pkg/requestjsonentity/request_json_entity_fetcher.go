package requestjsonentity

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/asaskevich/govalidator"
)

type Fetcher[T any] struct {
}

func (c *Fetcher[T]) Fetch(r io.Reader) (*T, error) {
	var entity *T
	dec := json.NewDecoder(r)
	err := dec.Decode(&entity)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	validated, err := govalidator.ValidateStruct(entity)
	if err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	if !validated {
		return nil, fmt.Errorf("validate request")
	}
	return entity, nil
}

func (c *Fetcher[T]) FetchSlice(r io.Reader) ([]T, error) {
	var list []T
	dec := json.NewDecoder(r)
	err := dec.Decode(&list)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	for _, item := range list {
		if _, err := govalidator.ValidateStruct(item); err != nil {
			return nil, fmt.Errorf("failed to validate request: %w", err)
		}
	}
	return list, nil
}
