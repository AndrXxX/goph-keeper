package hashgenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashGeneratorGenerate(t *testing.T) {
	tests := []struct {
		name string
		key  string
		data []byte
		want string
	}{
		{
			name: "Test with key 123 & data",
			key:  "123",
			data: []byte("data"),
			want: "3132333a6eb0790f39ac87c94f3856b2dd2c5d110e6811602261a9a923d3bb23adc8b7",
		},
		{
			name: "Test with key 123 & empty data",
			key:  "123",
			data: []byte(""),
			want: "313233e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "Test with empty key & empty data",
			key:  "",
			data: []byte(""),
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "Test with empty key & data",
			key:  "",
			data: []byte("data"),
			want: "3a6eb0790f39ac87c94f3856b2dd2c5d110e6811602261a9a923d3bb23adc8b7",
		},
	}
	f := Factory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := f.SHA256(tt.key)
			got := g.Generate(tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}
