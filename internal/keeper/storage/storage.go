package storage

import (
	"context"
	"time"

	"github.com/Julia-ivv/info-keeper.git/internal/keeper/config"
)

// Customer интерфейс для работы с пользователем.
type Customer interface {
	RegUser(ctx context.Context, login string, pwd string) error
	AuthUser(ctx context.Context, login string, pwd string) error
}

// CardWorker интерфейс для работы с банковскими картами.
type CardWorker interface {
	AddCard(ctx context.Context, userLogin string, prompt []byte,
		number []byte, date []byte, code []byte, note []byte, timeStamp time.Time) (err error)
	GetUserCardsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (cards []Card, err error)
	GetCard(ctx context.Context, userLogin string, number []byte) (card Card, err error)
	ForceUpdateCard(ctx context.Context, userLogin string, prompt []byte,
		number []byte, date []byte, code []byte, note []byte, timeStamp time.Time) (err error)
}

// LoginPwdWorker интерфейс для работы с парами логин-пароль.
type LoginPwdWorker interface {
	AddLoginPwd(ctx context.Context, userLogin string, prompt []byte,
		login []byte, pwd []byte, note []byte, timeStamp time.Time) (err error)
	GetUserLoginsPwdsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (loginsPwds []LoginPwd, err error)
	GetLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte) (loginPwd LoginPwd, err error)
	ForceUpdateLoginPwd(ctx context.Context, userLogin string, prompt []byte,
		login []byte, pwd []byte, note []byte, timeStamp time.Time) (err error)
}

// TextDataWorker интерфейс для работы с текстовыми данными.
type TextDataWorker interface {
	AddTextRecord(ctx context.Context, userLogin string, prompt []byte,
		data []byte, note []byte, timeStamp time.Time) (err error)
	GetUserTextRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []TextRecord, err error)
	GetTextRecord(ctx context.Context, userLogin string, prompt []byte) (record TextRecord, err error)
	ForceUpdateTextRecord(ctx context.Context, userLogin string, prompt []byte,
		data []byte, note []byte, timeStamp time.Time) (err error)
}

// BinaryDataWorker интерфейс для работы с бинарными данными.
type BinaryDataWorker interface {
	AddBinaryRecord(ctx context.Context, userLogin string, prompt []byte,
		data []byte, note []byte, timeStamp time.Time) (err error)
	GetUserBinaryRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []BinaryRecord, err error)
	GetBinaryRecord(ctx context.Context, userLogin string, prompt []byte) (record BinaryRecord, err error)
	ForceUpdateBinaryRecord(ctx context.Context, userLogin string, prompt []byte,
		data []byte, note []byte, timeStamp time.Time) (err error)
}

// Repositorier интерфейс для работы с репозиторием.
type Repositorier interface {
	Close() error
	Customer
	CardWorker
	LoginPwdWorker
	TextDataWorker
	BinaryDataWorker
}

// NewStorage создает новый объект репозитория.
func NewStorage(cfg config.Flags) (Repositorier, error) {
	db, err := NewDBStorage(cfg.DBDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
