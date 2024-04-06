package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

const LengthSalt = 4

func hash(value, salt string) string {
	var s = append([]byte(value), []byte(salt)...)
	hash := sha256.Sum256(s)
	hashString := hex.EncodeToString(hash[:])

	return hashString
}
