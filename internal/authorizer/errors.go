package authorizer

import "fmt"

// TypeAuthErrors - тип ошибок авторизации.
type TypeAuthErrors string

const (
	// QeuryError - ошибка выполнения запроса.
	QeuryError TypeAuthErrors = "request completed with error"
	// InvalidHash - ошибка недействительный хеш.
	InvalidHash TypeAuthErrors = "invalid hash"
)

// AuthErr - структура для ошибок авторизации.
type AuthErr struct {
	ErrType TypeAuthErrors
	Err     error
}

// Error - реализация интерфейса error.
func (e *AuthErr) Error() string {
	return fmt.Sprintln(e.ErrType)
}

// NewAuthError создает новую ошибку авторизации.
func NewAuthError(t TypeAuthErrors, err error) error {
	return &AuthErr{
		ErrType: t,
		Err:     err,
	}
}
