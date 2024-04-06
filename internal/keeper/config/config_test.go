package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	flags := NewConfig()
	if assert.NotEmpty(t, flags) {
		assert.NotEmpty(t, flags.GRPC)
	}
}

func TestReadFromConf(t *testing.T) {
	c := Flags{
		ConfigFileName: "for_tests.json",
	}
	err := readFromConf(&c)
	assert.NoError(t, err)
}
