package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Julia-ivv/info-keeper.git/pkg/randomizer"
	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	dbHandle *sql.DB
}

func createTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT UNIQUE NOT NULL CHECK(login != ''),
			hash TEXT NOT NULL CHECK(hash != ''),
			salt TEXT NOT NULL CHECK(salt != ''),
			last_sync TEXT NOT NULL CHECK(last_sync != '')
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS logins (
			user_id INTEGER NOT NULL REFERENCES users (user_id),
			prompt BLOB NOT NULL,
			login BLOB NOT NULL,
			pwd BLOB NOT NULL,
			note BLOB,
			time_stamp TEXT NOT NULL CHECK(time_stamp != ''),
			PRIMARY KEY(user_id, login, prompt)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS cards (
			user_id INTEGER NOT NULL REFERENCES users (user_id),
			prompt BLOB NOT NULL,
			number BLOB NOT NULL,
			date BLOB NOT NULL,
			code BLOB NOT NULL,
			note BLOB,
			time_stamp TEXT NOT NULL CHECK(time_stamp != ''),
			PRIMARY KEY(user_id, number)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS text_data (
			user_id INTEGER NOT NULL REFERENCES users (user_id),
			prompt BLOB NOT NULL,
			data BLOB NOT NULL,
			note BLOB,
			time_stamp TEXT NOT NULL CHECK(time_stamp != ''),
			PRIMARY KEY(user_id, prompt)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS binary_data (
			user_id INTEGER NOT NULL REFERENCES users (user_id),
			prompt BLOB NOT NULL,
			data BLOB NOT NULL,
			note BLOB,
			time_stamp TEXT NOT NULL CHECK(time_stamp != ''),
			PRIMARY KEY(user_id, prompt)
		)`)
	if err != nil {
		return err
	}

	return nil
}

func NewSQLiteStorage(DBURI string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", DBURI)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	err = createTables(ctx, db)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{dbHandle: db}, nil
}

func (db *SQLiteStorage) Close() error {
	return db.dbHandle.Close()
}

func (db *SQLiteStorage) RegUser(ctx context.Context, login string, pwd string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	salt, err := randomizer.GenerateRandomString(LengthSalt)
	if err != nil {
		return err
	}
	result, err := db.dbHandle.ExecContext(ctx,
		"INSERT INTO users (login, hash, salt, last_sync) VALUES (?,?,?,?)",
		login, hash(pwd, salt), salt, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("expected to affect 1 row")
	}

	return nil
}

func (db *SQLiteStorage) AuthUser(ctx context.Context, login string, pwd string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		"SELECT hash, salt FROM users WHERE login=?", login)

	var dbHash, dbSalt string
	err := row.Scan(&dbHash, &dbSalt)
	if err != nil {
		return err
	}

	newHash := hash(pwd, dbSalt)
	if newHash != dbHash {
		return errors.New("invalid hash")
	}

	return nil
}

type Card struct {
	Prompt    []byte
	Number    []byte
	Date      []byte
	Code      []byte
	Note      []byte
	TimeStamp []byte
}

func (db *SQLiteStorage) GetUserCardsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (cards []Card, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, number, date, code, note, time_stamp
		FROM cards
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND time_stamp > ?`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompt, number, date, code, note, timeStamp []byte
	for rows.Next() {
		err = rows.Scan(&prompt, &number, &date, &code, &note, &timeStamp)
		if err != nil {
			return nil, err
		}

		cards = append(cards, Card{
			Prompt:    prompt,
			Number:    number,
			Date:      date,
			Code:      code,
			Note:      note,
			TimeStamp: timeStamp,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return cards, nil
}

type LoginPwd struct {
	Prompt    []byte
	Login     []byte
	Pwd       []byte
	Note      []byte
	TimeStamp []byte
}

func (db *SQLiteStorage) GetUserLoginsPwdsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (loginsPwds []LoginPwd, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, login, pwd, note, time_stamp
		FROM logins
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND time_stamp > ?`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompt, login, pwd, note, timeStamp []byte
	for rows.Next() {
		err = rows.Scan(&prompt, &login, &pwd, &note, &timeStamp)
		if err != nil {
			return nil, err
		}
		loginsPwds = append(loginsPwds, LoginPwd{
			Prompt:    prompt,
			Login:     login,
			Pwd:       pwd,
			Note:      note,
			TimeStamp: timeStamp,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return loginsPwds, nil
}

type TextRecord struct {
	Prompt    []byte
	Data      []byte
	Note      []byte
	TimeStamp []byte
}

func (db *SQLiteStorage) GetUserTextRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []TextRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, data, note, time_stamp
		FROM text_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND time_stamp > ?`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prompt, data, note, timeStamp []byte
		err = rows.Scan(&prompt, &data, &note, &timeStamp)
		if err != nil {
			return nil, err
		}
		records = append(records, TextRecord{
			Prompt:    prompt,
			Data:      data,
			Note:      note,
			TimeStamp: timeStamp,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
}

type BinaryRecord struct {
	Prompt    []byte
	Data      []byte
	Note      []byte
	TimeStamp []byte
}

func (db *SQLiteStorage) GetUserBinaryRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []BinaryRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, data, note, time_stamp
		FROM binary_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND time_stamp > ?`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prompt, data, note, timeStamp []byte
		err = rows.Scan(&prompt, &data, &note, &timeStamp)
		if err != nil {
			return nil, err
		}
		records = append(records, BinaryRecord{
			Prompt:    prompt,
			Data:      data,
			Note:      note,
			TimeStamp: timeStamp,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
}
