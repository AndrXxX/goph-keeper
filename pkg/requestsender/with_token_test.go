package requestsender

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

type testTokenProvider struct {
	t string
}

func (p *testTokenProvider) GetToken() string {
	return p.t
}

func TestWithToken(t *testing.T) {
	tests := []struct {
		name          string
		tokenProvider tokenProvider
		params        dto.ParamsDto
		want          map[string]string
		wantErr       bool
	}{
		{
			name:          "Test with empty token",
			tokenProvider: &testTokenProvider{},
			params:        dto.ParamsDto{Headers: make(map[string]string)},
			want:          map[string]string{},
			wantErr:       false,
		},
		{
			name:          "Test with token test",
			tokenProvider: &testTokenProvider{"test"},
			params:        dto.ParamsDto{Headers: make(map[string]string)},
			want:          map[string]string{"Authorization": "Bearer test"},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := WithToken(tt.tokenProvider)
			require.Equal(t, tt.wantErr, f(&tt.params) != nil)
			assert.Equal(t, tt.want, tt.params.Headers)
		})
	}
}
