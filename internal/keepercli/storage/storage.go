package storage

import (
	"context"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/config"
)

// Customer интерфейс для работы с пользователем.
type Customer interface {
	RegUser(ctx context.Context, login string, pwd string) error
	AuthUser(ctx context.Context, login string, pwd string) error
}

// Synchronizer интерфейс для выполнения синхронизации.
type Synchronizer interface {
	GetLastSyncTime(ctx context.Context, userLogin string) (lastSync string, err error)
	AddSyncData(ctx context.Context, userLogin string,
		cards []Card, logins []LoginPwd, texts []TextRecord, binarys []BinaryRecord) (err error)
	UpdateLastSyncTime(ctx context.Context, userLogin string, syncTime string) (err error)
}

// CardWorker интерфейс для работы с банковскими картами.
type CardWorker interface {
	GetUserCardsAfterTime(ctx context.Context, userLogin string, afterTime string) (cards []Card, err error)
	AddCard(ctx context.Context, userLogin string, prompt []byte, number []byte, date []byte,
		code []byte, note []byte, timeStamp string) (err error)
	GetCard(ctx context.Context, userLogin string, number []byte) (card Card, err error)
	UpdateCard(ctx context.Context, userLogin string, prompt []byte, number []byte, date []byte,
		code []byte, note []byte, timeStamp string) (err error)
}

// LoginPwdWorker интерфейс для работы с парами логин-пароль.
type LoginPwdWorker interface {
	GetUserLoginsPwdsAfterTime(ctx context.Context, userLogin string, afterTime string) (loginsPwds []LoginPwd, err error)
	AddLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte, pwd []byte,
		note []byte, timeStamp string) (err error)
	GetLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte) (loginPwd LoginPwd, err error)
	UpdateLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte,
		pwd []byte, note []byte, timeStamp string) (err error)
}

// TextDataWorker интерфейс для работы с текстовыми данными.
type TextDataWorker interface {
	GetUserTextRecordsAfterTime(ctx context.Context, userLogin string, afterTime string) (records []TextRecord, err error)
	AddTextRecord(ctx context.Context, userLogin string, prompt []byte, data []byte,
		note []byte, timeStamp string) (err error)
	GetTextRecord(ctx context.Context, userLogin string, prompt []byte) (record TextRecord, err error)
	UpdateTextRecord(ctx context.Context, userLogin string, prompt []byte, data []byte,
		note []byte, timeStamp string) (err error)
}

// BinaryDataWorker интерфейс для работы с бинарными данными.
type BinaryDataWorker interface {
	GetUserBinaryRecordsAfterTime(ctx context.Context, userLogin string, afterTime string) (records []BinaryRecord, err error)
	AddBinaryRecord(ctx context.Context, userLogin string, prompt []byte, data []byte,
		note []byte, timeStamp string) (err error)
	GetBinaryRecord(ctx context.Context, userLogin string, prompt []byte) (record BinaryRecord, err error)
	UpdateBinaryRecord(ctx context.Context, userLogin string, prompt []byte, data []byte,
		note []byte, timeStamp string) (err error)
}

// Repositorier интерфейс для работы с репозиторием.
type Repositorier interface {
	Close() error
	Customer
	Synchronizer
	CardWorker
	LoginPwdWorker
	TextDataWorker
	BinaryDataWorker
}

// NewStorage создает новый объект репозитория.
func NewStorage(cfg config.Flags) (Repositorier, error) {
	db, err := NewSQLiteStorage(cfg.DBURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}
