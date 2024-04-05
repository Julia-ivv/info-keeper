package randomizer

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(length int) (string, error) {
	b, err := GenerateRandomBytes(length)
	return base64.RawURLEncoding.EncodeToString(b), err
}
