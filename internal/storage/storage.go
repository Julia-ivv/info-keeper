package storage

import "github.com/Julia-ivv/info-keeper.git/internal/config"

// type Customer interface {
// 	RegUser(ctx context.Context, regData RequestRegData) error
// 	AuthUser(ctx context.Context, authData RequestAuthData) error
// }

type Repositorier interface {
	Close() error
	// Customer
}

func NewStorage(cfg config.Flags) (Repositorier, error) {
	db, err := NewDBStorage(cfg.DBDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
