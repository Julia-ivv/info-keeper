package storage

import "fmt"

type TypeStorErrors string

const (
	// AlreadyExists - данные уже существуют в БД.
	AlreadyExists TypeStorErrors = "this data already exists in db"
	// ExistsDataNewerVersion - в БД существуют более новые данные.
	ExistsDataNewerVersion TypeStorErrors = "there is data with a newer version"
	// EmptyValues - пустое обязательное поле.
	EmptyValues TypeStorErrors = "empty required fields"
	// NullValues - значение обязательного полея null.
	NullValues TypeStorErrors = "required values are null"
	// EncryptionError - ошибка шифрования данных.
	EncryptionError TypeStorErrors = "encryption error"
	// DecryptionError - ошибка дешифрования данных.
	DecryptionError TypeStorErrors = "decryption error"
	// TypeStorErrors - запрашиваемые данные не найдены.
	EmptyResult TypeStorErrors = "the requested data not found"
)

// StorErr тип ошибок репозитория.
type StorErr struct {
	ErrType TypeStorErrors
	Err     error
}

// Error - реализация интерфейса error.
func (e *StorErr) Error() string {
	return fmt.Sprintln(e.ErrType)
}

// NewStorError создает новый объект ошибки репозитория.
func NewStorError(t TypeStorErrors, err error) error {
	return &StorErr{
		ErrType: t,
		Err:     err,
	}
}
