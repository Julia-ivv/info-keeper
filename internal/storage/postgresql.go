package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dbHandle *sql.DB
}

func NewDBStorage(DBURI string) (*DBStorage, error) {
	db, err := sql.Open("pgx", DBURI)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users (
			user_id serial, 
			login text UNIQUE NOT NULL, 
			hash text NOT NULL,
			salt text NOT NULL, 
			PRIMARY KEY(user_id)
		)`)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS logins (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt text,
			login text,
			pwd text,
			metadata text,
			PRIMARY KEY(user_id, login, prompt)
		)`)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS cards (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt text,
			number text,
			date text,
			code text,
			metadata text,
			PRIMARY KEY(user_id, number)
		)`)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS text_data (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt text,
			data text,
			metadata text,
			PRIMARY KEY(user_id, prompt)
		)`)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS binary_data (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt text,
			data bytea,
			metadata text,
			PRIMARY KEY(user_id, prompt)
		)`)
	if err != nil {
		return nil, err
	}

	return &DBStorage{dbHandle: db}, nil
}

func (db *DBStorage) Close() error {
	return db.dbHandle.Close()
}
