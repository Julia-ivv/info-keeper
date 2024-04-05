package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	l := NewLogger()
	assert.NotEmpty(t, l)
}
