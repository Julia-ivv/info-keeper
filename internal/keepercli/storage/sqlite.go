package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/Julia-ivv/info-keeper.git/pkg/randomizer"
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

// NewSQLiteStorage создает новый объект для работы с БД.
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

// Close закрывает БД.
func (db *SQLiteStorage) Close() error {
	return db.dbHandle.Close()
}

// RegUser регистрирует и аутентифицирует пользователя.
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

// AuthUser аутентифицирует пользователя.
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

// Card хранит информацию о банковской карте.
type Card struct {
	Prompt    []byte
	Number    []byte
	Date      []byte
	Code      []byte
	Note      []byte
	TimeStamp string
}

// GetUserCardsAfterTime получает информацию о банковских картах пользователя,
// введенную или измененную после указанного времени.
func (db *SQLiteStorage) GetUserCardsAfterTime(ctx context.Context, userLogin string,
	afterTime string) (cards []Card, err error) {
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

	var prompt, number, date, code, note []byte
	var timeStamp string
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

// GetCard получает информацию о банковской карте пользователя.
func (db *SQLiteStorage) GetCard(ctx context.Context, userLogin string, number []byte) (card Card, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT prompt, date, code, note, time_stamp
		FROM cards
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND number = ?`, userLogin, number)

	var prompt, date, code, note []byte
	var timeStamp string
	err = row.Scan(&prompt, &date, &code, &note, &timeStamp)
	if err != nil {
		return Card{}, err
	}

	return Card{
		Prompt:    prompt,
		Number:    number,
		Date:      date,
		Code:      code,
		Note:      note,
		TimeStamp: timeStamp,
	}, nil
}

// LoginPwd хранит информацию о паре логин-пароль.
type LoginPwd struct {
	Prompt    []byte
	Login     []byte
	Pwd       []byte
	Note      []byte
	TimeStamp string
}

// GetUserLoginsPwdsAfterTime получает информацию и парах логин-пароль,
// введенную или измененную после указанного времени.
func (db *SQLiteStorage) GetUserLoginsPwdsAfterTime(ctx context.Context, userLogin string,
	afterTime string) (loginsPwds []LoginPwd, err error) {
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

	var prompt, login, pwd, note []byte
	var timeStamp string
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

// GetLoginPwd получает данные о паре логин-пароль.
func (db *SQLiteStorage) GetLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte) (loginPwd LoginPwd, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT pwd, note, time_stamp
		FROM logins
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND prompt = ? AND login = ?`, userLogin, prompt, login)

	var pwd, note []byte
	var timeStamp string
	err = row.Scan(&pwd, &note, &timeStamp)
	if err != nil {
		return LoginPwd{}, err
	}

	return LoginPwd{
		Prompt:    prompt,
		Login:     login,
		Pwd:       pwd,
		Note:      note,
		TimeStamp: timeStamp,
	}, nil
}

// TextRecord хранит текстовую информацию.
type TextRecord struct {
	Prompt    []byte
	Data      []byte
	Note      []byte
	TimeStamp string
}

// GetUserTextRecordsAfterTime получает всю текстовую информацию пользователя,
// добавленную или измененную после указанного времени.
func (db *SQLiteStorage) GetUserTextRecordsAfterTime(ctx context.Context, userLogin string,
	afterTime string) (records []TextRecord, err error) {
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
		var prompt, data, note []byte
		var timeStamp string
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

// GetTextRecord получает текстовые данные.
func (db *SQLiteStorage) GetTextRecord(ctx context.Context, userLogin string, prompt []byte) (record TextRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT data, note, time_stamp
		FROM text_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND prompt = ?`, userLogin, prompt)

	var data, note []byte
	var timeStamp string
	err = row.Scan(&data, &note, &timeStamp)
	if err != nil {
		return TextRecord{}, err
	}

	return TextRecord{
		Prompt:    prompt,
		Data:      data,
		Note:      note,
		TimeStamp: timeStamp,
	}, nil
}

// BinaryRecord хранит бинарные данные.
type BinaryRecord struct {
	Prompt    []byte
	Data      []byte
	Note      []byte
	TimeStamp string
}

// GetUserBinaryRecordsAfterTime получает бинарную информацию пользователя,
// добавленную или измененную после указанного времени.
func (db *SQLiteStorage) GetUserBinaryRecordsAfterTime(ctx context.Context, userLogin string,
	afterTime string) (records []BinaryRecord, err error) {
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
		var prompt, data, note []byte
		var timeStamp string
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

// GetBinaryRecord получает бинарную информацию.
func (db *SQLiteStorage) GetBinaryRecord(ctx context.Context, userLogin string, prompt []byte) (record BinaryRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT data, note, time_stamp
		FROM binary_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND prompt = ?`, userLogin, prompt)

	var data, note []byte
	var timeStamp string
	err = row.Scan(&data, &note, &timeStamp)
	if err != nil {
		return BinaryRecord{}, err
	}

	return BinaryRecord{
		Prompt:    prompt,
		Data:      data,
		Note:      note,
		TimeStamp: timeStamp,
	}, nil
}

// GetLastSyncTime получает время последней синхронизации.
func (db *SQLiteStorage) GetLastSyncTime(ctx context.Context, userLogin string) (lastSync string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT last_sync
		FROM users
		WHERE login = ?`, userLogin)

	err = row.Scan(&lastSync)
	if err != nil {
		return "", err
	}

	return lastSync, nil
}

// UpdateLastSyncTime обновляет время последней синхронизации.
func (db *SQLiteStorage) UpdateLastSyncTime(ctx context.Context, userLogin string, syncTime string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE users 
		SET last_sync = ?
		WHERE login = ?`, syncTime, userLogin)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// AddCard добавляет информацию о банковской карте.
func (db *SQLiteStorage) AddCard(ctx context.Context, userLogin string, prompt []byte, number []byte, date []byte,
	code []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO cards (user_id , prompt, number, date, code, note, time_stamp) 
				VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?,?,?)`,
		userLogin, prompt, number, date, code, note, timeStamp)
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

// AddLoginPwd добавляет информацию о паре логин-пароль.
func (db *SQLiteStorage) AddLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte,
	pwd []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO logins (user_id , prompt, login, pwd, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?,?)`,
		userLogin, prompt, login, pwd, note, timeStamp)
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

// AddTextRecord добавляет текстовую информацию.
func (db *SQLiteStorage) AddTextRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO text_data (user_id , prompt, data, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?)`,
		userLogin, prompt, data, note, timeStamp)
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

// AddBinaryRecord добавляет бинарную информацию.
func (db *SQLiteStorage) AddBinaryRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO binary_data (user_id , prompt, data, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?)`,
		userLogin, prompt, data, note, timeStamp)
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

// AddSyncData добавляет новые данные, полученные от сервера при синхронизации.
func (db *SQLiteStorage) AddSyncData(ctx context.Context, userLogin string,
	cards []Card, logins []LoginPwd, texts []TextRecord, binarys []BinaryRecord) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := db.dbHandle.Begin()
	if err != nil {
		return err
	}

	for _, v := range cards {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO cards (user_id , prompt, number, date, code, note, time_stamp) 
					VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?,?,?)`,
			userLogin, v.Prompt, v.Number, v.Date, v.Code, v.Note, v.TimeStamp)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if rows != 1 {
			tx.Rollback()
			return fmt.Errorf("expected to affect 1 row, affected %d", rows)
		}
	}

	for _, v := range logins {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO logins (user_id , prompt, login, pwd, note, time_stamp) 
			VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?,?)`,
			userLogin, v.Prompt, v.Login, v.Pwd, v.Note, v.TimeStamp)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if rows != 1 {
			tx.Rollback()
			return fmt.Errorf("expected to affect 1 row, affected %d", rows)
		}
	}

	for _, v := range texts {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO text_data (user_id , prompt, data, note, time_stamp) 
			VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?)`,
			userLogin, v.Prompt, v.Data, v.Note, v.TimeStamp)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if rows != 1 {
			tx.Rollback()
			return fmt.Errorf("expected to affect 1 row, affected %d", rows)
		}
	}

	for _, v := range binarys {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO binary_data (user_id , prompt, data, note, time_stamp) 
			VALUES ((SELECT user_id FROM users WHERE login = ?),?,?,?,?)`,
			userLogin, v.Prompt, v.Data, v.Note, v.TimeStamp)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if rows != 1 {
			tx.Rollback()
			return fmt.Errorf("expected to affect 1 row, affected %d", rows)
		}
	}

	return tx.Commit()
}

// UpdateCard обновляет информацию о банковской карте.
func (db *SQLiteStorage) UpdateCard(ctx context.Context, userLogin string, prompt []byte,
	number []byte, date []byte, code []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE cards 
		SET prompt = ?, date = ?, code = ?, note = ?, time_stamp = ?
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND number = ?`,
		prompt, date, code, note, timeStamp, userLogin, number)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// UpdateLoginPwd обновляет информацию о паре логин-пароль.
func (db *SQLiteStorage) UpdateLoginPwd(ctx context.Context, userLogin string, prompt []byte,
	login []byte, pwd []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE logins 
		SET pwd = ?, note = ?, time_stamp = ?
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND prompt = ?
		AND login = ?`,
		pwd, note, timeStamp, userLogin, prompt, login)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// UpdateTextRecord обновляет текстовую информацию.
func (db *SQLiteStorage) UpdateTextRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE text_data 
		SET data = ?, note = ?, time_stamp = ?
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND prompt = ?`,
		data, note, timeStamp, userLogin, prompt)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// UpdateBinaryRecord обновляет бинарные данные.
func (db *SQLiteStorage) UpdateBinaryRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE binary_data 
		SET data = ?, note = ?, time_stamp = ?
		WHERE user_id = (SELECT user_id FROM users WHERE login = ?)
		AND prompt = ?`,
		data, note, timeStamp, userLogin, prompt)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}
