package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	v := "data"
	s := "salt"
	h := hash(v, s)
	assert.NotEmpty(t, h)
}
