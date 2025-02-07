package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotEmpty(t, config.Host)
	assert.NotEmpty(t, config.LogLevel)
}
