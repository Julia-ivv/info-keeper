package storage

import (
	"context"
	"time"

	"github.com/Julia-ivv/info-keeper.git/internal/config"
)

type Customer interface {
	RegUser(ctx context.Context, login string, pwd string) error
	AuthUser(ctx context.Context, login string, pwd string) error
}

type CardWorker interface {
	AddCard(ctx context.Context, userLogin string, prompt string,
		number string, date string, code string, note string, timeStamp time.Time) (err error)
	GetUserCardsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (cards []Card, err error)
	GetCard(ctx context.Context, userLogin string, number string) (card Card, err error)
	ForceUpdateCard(ctx context.Context, userLogin string, prompt string,
		number string, date string, code string, note string, timeStamp time.Time) (err error)
}

type LoginPwdWorker interface {
	AddLoginPwd(ctx context.Context, userLogin string, prompt string,
		login string, pwd string, note string, timeStamp time.Time) (err error)
	GetUserLoginsPwdsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (loginsPwds []LoginPwd, err error)
	GetLoginPwd(ctx context.Context, userLogin string, prompt string, login string) (loginPwd LoginPwd, err error)
	ForceUpdateLoginPwd(ctx context.Context, userLogin string, prompt string,
		login string, pwd string, note string, timeStamp time.Time) (err error)
}

type TextDataWorker interface {
	AddTextRecord(ctx context.Context, userLogin string, prompt string,
		data string, note string, timeStamp time.Time) (err error)
	GetUserTextRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []TextRecord, err error)
	GetTextRecord(ctx context.Context, userLogin string, prompt string) (record TextRecord, err error)
	ForceUpdateTextRecord(ctx context.Context, userLogin string, prompt string,
		data string, note string, timeStamp time.Time) (err error)
}

type BinaryDataWorker interface {
	AddBinaryRecord(ctx context.Context, userLogin string, prompt string,
		data []byte, note string, timeStamp time.Time) (err error)
	GetUserBinaryRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []BinaryRecord, err error)
	GetBinaryRecord(ctx context.Context, userLogin string, prompt string) (record BinaryRecord, err error)
	ForceUpdateBinaryRecord(ctx context.Context, userLogin string, prompt string,
		data []byte, note string, timeStamp time.Time) (err error)
}

type Repositorier interface {
	Close() error
	Customer
	CardWorker
	LoginPwdWorker
	TextDataWorker
	BinaryDataWorker
}

func NewStorage(cfg config.Flags) (Repositorier, error) {
	db, err := NewDBStorage(cfg.DBDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
