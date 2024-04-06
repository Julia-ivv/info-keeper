package storage

import "fmt"

type TypeStorErrors string

const (
	AlreadyExists          TypeStorErrors = "this data already exists in db"
	ExistsDataNewerVersion TypeStorErrors = "there is data with a newer version"
	EmptyValues            TypeStorErrors = "empty required fields"
	NullValues             TypeStorErrors = "required values are null"
	EncryptionError        TypeStorErrors = "encryption error"
	DecryptionError        TypeStorErrors = "decryption error"
	EmptyResult            TypeStorErrors = "the requested data not found"
)

type StorErr struct {
	ErrType TypeStorErrors
	Err     error
}

func (e *StorErr) Error() string {
	return fmt.Sprintln(e.ErrType)
}

func NewStorError(t TypeStorErrors, err error) error {
	return &StorErr{
		ErrType: t,
		Err:     err,
	}
}
