package configprovider

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	Host string `valid:"minstringlength(3)"`
}

type tempParser struct {
	err  error
	host string
}

func (p tempParser) Parse(c *testConfig) error {
	c.Host = p.host
	return p.err
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		parsers []parser[testConfig]
		want    *configProvider[testConfig]
	}{
		{
			name:    "Test with empty parsers",
			parsers: []parser[testConfig]{},
			want:    &configProvider[testConfig]{parsers: []parser[testConfig]{}},
		},
		{
			name:    "Test with temp parser",
			parsers: []parser[testConfig]{tempParser{}},
			want:    &configProvider[testConfig]{parsers: []parser[testConfig]{tempParser{}}},
		},
		{
			name:    "Test with two temp parsers",
			parsers: []parser[testConfig]{tempParser{}, tempParser{}},
			want:    &configProvider[testConfig]{parsers: []parser[testConfig]{tempParser{}, tempParser{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(nil, tt.parsers...))
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		parsers []parser[testConfig]
		wantErr bool
	}{
		{
			name:    "Test with err parser",
			parsers: []parser[testConfig]{tempParser{err: errors.New("err")}},
			wantErr: true,
		},
		{
			name:    "Test with no err parser",
			parsers: []parser[testConfig]{tempParser{}},
			wantErr: false,
		},
		{
			name:    "Test with validate err",
			parsers: []parser[testConfig]{tempParser{host: "-"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := New(&testConfig{}, tt.parsers...)
			_, err := provider.Fetch()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
