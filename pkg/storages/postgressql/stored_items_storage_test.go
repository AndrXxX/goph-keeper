package postgressql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vingarcia/ksql"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
	"github.com/AndrXxX/goph-keeper/pkg/utils"
)

func Test_storedItemsStorage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		db      ksql.Provider
		m       *models.StoredItem
		want    *models.StoredItem
		wantErr bool
	}{
		{
			name: "Test with error",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
				return &p
			}(),
			m:       &models.StoredItem{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test OK",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &p
			}(),
			m:       &models.StoredItem{UserID: 10},
			want:    &models.StoredItem{UserID: 10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory{DB: tt.db}
			s := f.StoredItemsStorage()
			got, err := s.Insert(context.Background(), tt.m)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_storedItemsStorage_Query(t *testing.T) {
	tests := []struct {
		name    string
		db      ksql.Provider
		m       *models.StoredItem
		want    []models.StoredItem
		wantErr bool
	}{
		{
			name: "Test with error",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
				return &p
			}(),
			m:       &models.StoredItem{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test OK",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &p
			}(),
			m:       &models.StoredItem{UpdatedAt: utils.Pointer(time.Now())},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory{DB: tt.db}
			s := f.StoredItemsStorage()
			got, err := s.Query(context.Background(), tt.m)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_storedItemsStorage_QueryOneById(t *testing.T) {
	tests := []struct {
		name    string
		db      ksql.Provider
		id      uuid.UUID
		want    *models.StoredItem
		wantErr bool
	}{
		{
			name: "Test with error",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("QueryOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
				return &p
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test OK",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("QueryOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &p
			}(),
			id:      uuid.UUID{},
			want:    &models.StoredItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory{DB: tt.db}
			s := f.StoredItemsStorage()
			got, err := s.QueryOneById(context.Background(), tt.id)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_storedItemsStorage_Update(t *testing.T) {
	tests := []struct {
		name    string
		db      ksql.Provider
		m       *models.StoredItem
		want    *models.StoredItem
		wantErr bool
	}{
		{
			name: "Test with error",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Patch", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
				return &p
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test OK",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Patch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return &p
			}(),
			m:       &models.StoredItem{},
			want:    &models.StoredItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory{DB: tt.db}
			s := f.StoredItemsStorage()
			got, err := s.Update(context.Background(), tt.m)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
