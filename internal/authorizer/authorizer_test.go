package authorizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildToken(t *testing.T) {
	l := "login"
	p := "pwd"
	t.Run("test for build token", func(t *testing.T) {
		tokenStr, err := BuildToken(l, p, "key")
		assert.NotEmpty(t, tokenStr)
		assert.NoError(t, err)
	})
}

func TestGetUserIDFromToken(t *testing.T) {
	t.Run("test with random string", func(t *testing.T) {
		l, err := GetUserDataFromToken("some_string", "key")
		assert.Empty(t, l)
		assert.Error(t, err)
	})
}
