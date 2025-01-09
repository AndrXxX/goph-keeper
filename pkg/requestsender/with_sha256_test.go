package requestsender

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/goph-keeper/pkg/requestsender/dto"
)

type testHg struct {
	mock.Mock
}

func (g *testHg) Generate(data []byte) string {
	args := g.Called(data)
	return args.String(0)
}

func TestWithSHA256(t *testing.T) {
	tests := []struct {
		name    string
		hg      hashGenerator
		key     string
		params  dto.ParamsDto
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Test with error on read data",
			params: func() dto.ParamsDto {
				r := readerMock{}
				r.On("Read", mock.Anything).Return(0, errors.New("error on read"))
				return dto.ParamsDto{Buf: &r, Headers: make(map[string]string)}
			}(),
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name: "Test with hashed header",
			hg: func() hashGenerator {
				hg := testHg{}
				hg.On("Generate", mock.Anything).Return("hashedResult")
				return &hg
			}(),
			params:  dto.ParamsDto{Buf: bytes.NewReader([]byte("test")), Headers: make(map[string]string)},
			want:    map[string]string{"HashSHA256": "hashedResult"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := WithSHA256(tt.hg)
			require.Equal(t, tt.wantErr, f(&tt.params) != nil)
			assert.Equal(t, tt.want, tt.params.Headers)
		})
	}
}
