package cryptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNonce(t *testing.T) {
	n, aesgsm, err := generateNonce()
	assert.NoError(t, err)
	assert.NotEmpty(t, aesgsm)
	assert.NotEmpty(t, n)
}

func TestEncryptsString(t *testing.T) {
	b, err := EncryptsString("some data")
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}

func TestEncryptsByte(t *testing.T) {
	b, err := EncryptsByte([]byte{46, 85})
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}

func TestDecrypts(t *testing.T) {
	str := "some string"
	d, e := EncryptsString(str)
	if assert.NoError(t, e) {
		s, err := Decrypts(d)
		assert.NoError(t, err)
		assert.Equal(t, str, s)
	}
}

func TestDecryptsByte(t *testing.T) {
	bs := []byte{45, 45, 45}
	d, e := EncryptsByte(bs)
	if assert.NoError(t, e) {
		s, err := DecryptsByte(d)
		assert.NoError(t, err)
		assert.Equal(t, bs, s)
	}
}
