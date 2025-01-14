package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var id1 = uuid.New()

func TestStoredItem(t *testing.T) {
	tests := []struct {
		name         string
		entity       StoredItem
		wantId       uuid.UUID
		wantIsStored bool
	}{
		{
			name:         "Test with id1",
			entity:       StoredItem{ID: id1},
			wantId:       id1,
			wantIsStored: true,
		},
		{
			name:         "Test with empty id",
			entity:       StoredItem{},
			wantIsStored: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantId, tt.entity.GetID())
			assert.Equal(t, tt.wantIsStored, tt.entity.IsStored())
		})
	}
}
