package storage

import (
	"context"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/config"
)

type Customer interface {
	RegUser(ctx context.Context, login string, pwd string) error
	AuthUser(ctx context.Context, login string, pwd string) error
}

type Repositorier interface {
	Close() error
	Customer
}

func NewStorage(cfg config.Flags) (Repositorier, error) {
	db, err := NewSQLiteStorage(cfg.DBURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}
