package certgenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenCert(t *testing.T) {
	cFile, pFile, err := GenCert(1024)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, cFile)
		assert.NotEmpty(t, pFile)
	}
}
