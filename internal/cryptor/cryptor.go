package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

const (
	password = "x35k9fhrds45hf"
)

func generateNonce() (nonce []byte, aesgcm cipher.AEAD, err error) {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}

	aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		return nil, nil, err
	}

	return key[len(key)-aesgcm.NonceSize():], aesgcm, nil
}

func EncryptsString(data string) (result []byte, err error) {
	nonce, aesgcm, err := generateNonce()
	if err != nil {
		return nil, err
	}

	res := aesgcm.Seal(nil, nonce, []byte(data), nil)
	return res, nil
}

func EncryptsByte(data []byte) (result []byte, err error) {
	nonce, aesgcm, err := generateNonce()
	if err != nil {
		return nil, err
	}

	res := aesgcm.Seal(nil, nonce, data, nil)
	return res, nil
}

func Decrypts(data []byte) (result string, err error) {
	nonce, aesgcm, err := generateNonce()
	if err != nil {
		return "", err
	}

	res, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func DecryptsByte(data []byte) (result []byte, err error) {
	nonce, aesgcm, err := generateNonce()
	if err != nil {
		return nil, err
	}

	res, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
