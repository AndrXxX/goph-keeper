package envparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
)

func Test_parseEnv(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		env     map[string]string
		want    *config.Config
		wantErr bool
	}{
		{
			name:    "Empty env",
			config:  &config.Config{Host: "host"},
			env:     map[string]string{},
			want:    &config.Config{Host: "host"},
			wantErr: false,
		},
		{
			name:    "HOST=new-host",
			config:  &config.Config{Host: "host"},
			env:     map[string]string{"HOST": "new-host"},
			want:    &config.Config{Host: "new-host"},
			wantErr: false,
		},
		{
			name:    "DATABASE_URI=test",
			config:  &config.Config{DatabaseURI: ""},
			env:     map[string]string{"DATABASE_URI": "test"},
			want:    &config.Config{DatabaseURI: "test"},
			wantErr: false,
		},
		{
			name:    "AUTH_SECRET_KEY=abc",
			config:  &config.Config{AuthKey: ""},
			env:     map[string]string{"AUTH_SECRET_KEY": "abc"},
			want:    &config.Config{AuthKey: "abc"},
			wantErr: false,
		},
		{
			name:    "AUTH_SECRET_KEY_EXPIRED=abc",
			config:  &config.Config{},
			env:     map[string]string{"AUTH_SECRET_KEY_EXPIRED": "abc"},
			want:    &config.Config{},
			wantErr: true,
		},
	}
	parser := Parser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			err := parser.Parse(tt.config)
			assert.Equal(t, tt.want, tt.config)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
