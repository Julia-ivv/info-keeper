package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorError(t *testing.T) {
	e := NewStorError(AlreadyExists, errors.New("some error"))
	assert.Error(t, e)
	assert.Equal(t, &StorErr{
		ErrType: AlreadyExists,
		Err:     errors.New("some error"),
	}, e)
}
