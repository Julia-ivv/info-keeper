package randomizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes(t *testing.T) {
	l := 5
	b, err := GenerateRandomBytes(l)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}

func TestGenerateRandomString(t *testing.T) {
	l := 8
	s, err := GenerateRandomString(l)
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}
