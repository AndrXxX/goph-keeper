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
	if _, err := govalidator.ValidateStruct(entity); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}
	return entity, nil
}
