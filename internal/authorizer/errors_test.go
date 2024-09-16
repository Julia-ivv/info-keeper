package authorizer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenError(t *testing.T) {
	t.Run("test create error struct", func(t *testing.T) {
		err := NewAuthError(QeuryError, errors.New("error"))
		assert.NotEmpty(t, err)
		assert.Equal(t, &AuthErr{
			ErrType: QeuryError,
			Err:     errors.New("error"),
		}, err)
	})
}
