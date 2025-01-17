package postgressql

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vingarcia/ksql"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

func Test_usersStorage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		db      ksql.Provider
		m       *models.User
		want    *models.User
		wantErr bool
	}{
		{
			name: "Test with error",
			db: func() ksql.Provider {
				p := ksqlProvider{}
				p.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
				return &p
			}(),
			m:       &models.User{},
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
			m:       &models.User{Login: "test"},
			want:    &models.User{Login: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory{DB: tt.db}
			s := f.UsersStorage()
			got, err := s.Insert(context.Background(), tt.m)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_usersStorage_QueryOne(t *testing.T) {
	tests := []struct {
		name    string
		db      ksql.Provider
		login   string
		want    *models.User
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
			login:   "test",
			want:    &models.User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory{DB: tt.db}
			s := f.UsersStorage()
			got, err := s.QueryOne(context.Background(), tt.login)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
