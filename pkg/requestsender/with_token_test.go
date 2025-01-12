package requestsender

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

func TestWithToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		params  dto.ParamsDto
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "Test with empty token",
			token:   "",
			params:  dto.ParamsDto{Headers: make(map[string]string)},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "Test with token test",
			token:   "test",
			params:  dto.ParamsDto{Headers: make(map[string]string)},
			want:    map[string]string{"Authorization": "Bearer test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := WithToken(tt.token)
			require.Equal(t, tt.wantErr, f(&tt.params) != nil)
			assert.Equal(t, tt.want, tt.params.Headers)
		})
	}
}
