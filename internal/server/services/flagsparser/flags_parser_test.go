package flagsparser

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
)

type testCase struct {
	name   string
	config *config.Config
	flags  []string
	want   *config.Config
}

func Test_parseFlags(t *testing.T) {
	tests := []testCase{
		{
			name:   "Empty flags",
			config: &config.Config{Host: "host"},
			flags:  []string{},
			want:   &config.Config{Host: "host"},
		},
		{
			name:   "-a=new-host",
			config: &config.Config{Host: "host"},
			flags:  []string{"-a", "new-host"},
			want:   &config.Config{Host: "new-host"},
		},
		{
			name:   "-ake=5",
			config: &config.Config{},
			flags:  []string{"-ake", "5"},
			want:   &config.Config{AuthKeyExpired: 5},
		},
		{
			name:   "-d=test",
			config: &config.Config{},
			flags:  []string{"-d", "test"},
			want:   &config.Config{DatabaseURI: "test"},
		},
	}
	for _, tt := range tests {
		run(t, tt)
	}
}

func run(t *testing.T, tt testCase) {
	t.Run(tt.name, func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = os.Args[:1]
		os.Args = append(os.Args[:1], tt.flags...)
		err := FlagsParser{}.Parse(tt.config)
		assert.Equal(t, tt.want, tt.config)
		assert.NoError(t, err)
	})
}
