package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/pkg/randomizer"
)

type DBStorage struct {
	dbHandle *sql.DB
}

func createTables(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users (
			user_id serial UNIQUE, 
			login text UNIQUE NOT NULL CHECK(login != ''), 
			hash text NOT NULL CHECK(hash != ''),
			salt text NOT NULL CHECK(salt != ''), 
			PRIMARY KEY(user_id)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS logins (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt bytea NOT NULL,
			login bytea NOT NULL,
			pwd bytea NOT NULL,
			note bytea,
			time_stamp timestamptz (0) NOT NULL,
			PRIMARY KEY(user_id, login, prompt)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS cards (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt bytea NOT NULL,
			number bytea NOT NULL,
			date bytea NOT NULL,
			code bytea NOT NULL,
			note bytea,
			time_stamp timestamptz (0) NOT NULL,
			PRIMARY KEY(user_id, number)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS text_data (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt bytea NOT NULL,
			data bytea NOT NULL,
			note bytea,
			time_stamp timestamptz (0) NOT NULL,
			PRIMARY KEY(user_id, prompt)
		)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS binary_data (
			user_id integer NOT NULL REFERENCES users(user_id),
			prompt bytea NOT NULL,
			data bytea NOT NULL,
			note bytea,
			time_stamp timestamptz (0) NOT NULL,
			PRIMARY KEY(user_id, prompt)
		)`)
	if err != nil {
		return err
	}

	return nil
}

// NewDBStorage создает объект для работы с БД.
func NewDBStorage(DBURI string) (*DBStorage, error) {
	db, err := sql.Open("pgx", DBURI)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = createTables(ctx, db)
	if err != nil {
		return nil, err
	}

	return &DBStorage{dbHandle: db}, nil
}

// Close закрывает БД.
func (db *DBStorage) Close() error {
	return db.dbHandle.Close()
}

// RegUser добавляет нового пользователя в БД.
func (db *DBStorage) RegUser(ctx context.Context, login string, pwd string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	salt, err := randomizer.GenerateRandomString(LengthSalt)
	if err != nil {
		return err
	}
	result, err := db.dbHandle.ExecContext(ctx,
		"INSERT INTO users (login, hash, salt) VALUES ($1, $2, $3)",
		login, hash(pwd, salt), salt)
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
func (db *DBStorage) AuthUser(ctx context.Context, login string, pwd string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		"SELECT hash, salt FROM users WHERE login=$1", login)

	var dbHash, dbSalt string
	err := row.Scan(&dbHash, &dbSalt)
	if err != nil {
		return authorizer.NewAuthError(authorizer.QeuryError, err)
	}

	newHash := hash(pwd, dbSalt)
	if newHash != dbHash {
		return authorizer.NewAuthError(authorizer.InvalidHash, errors.New("invalid hash"))
	}

	return nil
}

// AddCard добавляет информацию о банковской карте.
func (db *DBStorage) AddCard(ctx context.Context, userLogin string, prompt []byte,
	number []byte, date []byte, code []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO cards (user_id , prompt, number, date, code, note, time_stamp) 
			VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5, $6, $7)`,
		userLogin, prompt, number, date, code, note, timeStamp)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM cards 
				WHERE number = $1 AND 
				user_id = (SELECT user_id FROM users WHERE login = $2)`, number, userLogin)
			var tServer time.Time
			errScan := row.Scan(&tServer)
			if errScan != nil {
				return errScan
			}
			if tServer.After(timeStamp) {
				return NewStorError(ExistsDataNewerVersion, err)
			}
			result, err = db.dbHandle.ExecContext(ctx,
				`UPDATE cards 
				SET prompt = $1, date = $2, code = $3, note = $4, time_stamp = $5
				WHERE user_id = (SELECT user_id FROM users WHERE login = $6)
				AND number = $7`,
				prompt, date, code, note, timeStamp, userLogin, number)
			if err != nil {
				return err
			}
			rowUpd, errUpd := result.RowsAffected()
			if errUpd != nil {
				return errUpd
			}
			if rowUpd != 1 {
				return errors.New("expected to affect 1 row")
			}
			return nil
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			return NewStorError(NullValues, err)
		}
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
func (db *DBStorage) AddLoginPwd(ctx context.Context, userLogin string, prompt []byte,
	login []byte, pwd []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO logins (user_id , prompt, login, pwd, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5, $6)`,
		userLogin, prompt, login, pwd, note, timeStamp)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM logins 
				WHERE prompt = $1 
				AND login = $2 
				AND user_id = (SELECT user_id FROM users WHERE login = $3)`,
				prompt, login, userLogin)
			var tServer time.Time
			errScan := row.Scan(&tServer)
			if errScan != nil {
				return errScan
			}
			if tServer.After(timeStamp) {
				return NewStorError(ExistsDataNewerVersion, err)
			}
			result, err = db.dbHandle.ExecContext(ctx,
				`UPDATE logins 
				SET pwd = $1, note = $2, time_stamp = $3
				WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
				AND prompt = $5
				AND login = $6`,
				pwd, note, timeStamp, userLogin, prompt, login)
			if err != nil {
				return err
			}
			rowUpd, errUpd := result.RowsAffected()
			if errUpd != nil {
				return errUpd
			}
			if rowUpd != 1 {
				return errors.New("expected to affect 1 row")
			}
			return nil
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			return NewStorError(NullValues, err)
		}
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

// AddTextRecord раелизует добавление текстовой информации.
func (db *DBStorage) AddTextRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO text_data (user_id , prompt, data, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5)`,
		userLogin, prompt, data, note, timeStamp)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM text_data 
				WHERE prompt = $1 AND 
				user_id = (SELECT user_id FROM users WHERE login = $2)`, prompt, userLogin)
			var tServer time.Time
			errScan := row.Scan(&tServer)
			if errScan != nil {
				return errScan
			}
			if tServer.After(timeStamp) {
				return NewStorError(ExistsDataNewerVersion, err)
			}
			result, err = db.dbHandle.ExecContext(ctx,
				`UPDATE text_data 
				SET data = $1, note = $2, time_stamp = $3
				WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
				AND prompt = $5`,
				data, note, timeStamp, userLogin, prompt)
			if err != nil {
				return err
			}
			rowUpd, errUpd := result.RowsAffected()
			if errUpd != nil {
				return errUpd
			}
			if rowUpd != 1 {
				return errors.New("expected to affect 1 row")
			}
			return nil
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			return NewStorError(NullValues, err)
		}
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

// AddBinaryRecord реализует добавление бинарной информации в БД.
func (db *DBStorage) AddBinaryRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO binary_data (user_id , prompt, data, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5)`,
		userLogin, prompt, data, note, timeStamp)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM binary_data 
				WHERE prompt = $1 AND 
				user_id = (SELECT user_id FROM users WHERE login = $2)`, prompt, userLogin)
			var tServer time.Time
			errScan := row.Scan(&tServer)
			if errScan != nil {
				return errScan
			}
			if tServer.After(timeStamp) {
				return NewStorError(ExistsDataNewerVersion, err)
			}
			result, err = db.dbHandle.ExecContext(ctx,
				`UPDATE binary_data 
				SET data = $1, note = $2, time_stamp = $3
				WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
				AND prompt = $5`,
				data, note, timeStamp, userLogin, prompt)
			if err != nil {
				return err
			}
			rowUpd, errUpd := result.RowsAffected()
			if errUpd != nil {
				return errUpd
			}
			if rowUpd != 1 {
				return errors.New("expected to affect 1 row")
			}
			return nil
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			return NewStorError(NullValues, err)
		}
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

// Card хранит информацию о банковской карте.
type Card struct {
	Prompt    []byte
	Number    []byte
	Date      []byte
	Code      []byte
	Note      []byte
	TimeStamp time.Time
}

// GetCard получает информацию о банковской карте.
func (db *DBStorage) GetCard(ctx context.Context, userLogin string, number []byte) (card Card, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT prompt, date, code, note, time_stamp
		FROM cards
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND number = $2`, userLogin, number)

	var prompt, date, code, note []byte
	var timeStamp time.Time
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

// GetUserCardsAfterTime - получает все банковские карты пользователя,
// добавленные или измененные после указанного времени.
func (db *DBStorage) GetUserCardsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (cards []Card, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, number, date, code, note, time_stamp
		FROM cards
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND time_stamp > $2`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompt, number, date, code, note []byte
	var timeStamp time.Time
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

// LoginPwd хранит информацию о парах логин-пароль.
type LoginPwd struct {
	Prompt    []byte
	Login     []byte
	Pwd       []byte
	Note      []byte
	TimeStamp time.Time
}

// GetLoginPwd получает информацию о паре логин-пароль.
func (db *DBStorage) GetLoginPwd(ctx context.Context, userLogin string, prompt []byte, login []byte) (loginPwd LoginPwd, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT pwd, note, time_stamp
		FROM logins
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND prompt = $2 AND login = $3`, userLogin, prompt, login)

	var pwd, note []byte
	var timeStamp time.Time
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

// GetUserLoginsPwdsAfterTime получает информацию о парах логин-пароль пользователя,
// добавленных или измененных после указанного времени.
func (db *DBStorage) GetUserLoginsPwdsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (loginsPwds []LoginPwd, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, login, pwd, note, time_stamp
		FROM logins
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND time_stamp > $2`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompt, login, pwd, note []byte
	var timeStamp time.Time
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

// TextRecord хранит текстовую информацию.
type TextRecord struct {
	Prompt    []byte
	Data      []byte
	Note      []byte
	TimeStamp time.Time
}

// GetTextRecord получает текстовую информацию.
func (db *DBStorage) GetTextRecord(ctx context.Context, userLogin string, prompt []byte) (record TextRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT data, note, time_stamp
		FROM text_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND prompt = $2`, userLogin, prompt)

	var data, note []byte
	var timeStamp time.Time
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

// GetUserTextRecordsAfterTime получает все текстовые данные пользователя,
// добавленные или измененнные после указанного времени.
func (db *DBStorage) GetUserTextRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []TextRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, data, note, time_stamp
		FROM text_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND time_stamp > $2`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prompt, data, note []byte
		var timeStamp time.Time
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

// BinaryRecord хранит бинарные данные.
type BinaryRecord struct {
	Prompt    []byte
	Data      []byte
	Note      []byte
	TimeStamp time.Time
}

// GetBinaryRecord получает бинарные данные.
func (db *DBStorage) GetBinaryRecord(ctx context.Context, userLogin string, prompt []byte) (record BinaryRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT data, note, time_stamp
		FROM binary_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND prompt = $2`, userLogin, prompt)

	var data, note []byte
	var timeStamp time.Time
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

// GetUserBinaryRecordsAfterTime получает все бинарные данные пользователя,
// добавленные или измененные после указанного времени.
func (db *DBStorage) GetUserBinaryRecordsAfterTime(ctx context.Context, userLogin string, afterTime time.Time) (records []BinaryRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.dbHandle.QueryContext(ctx,
		`SELECT prompt, data, note, time_stamp
		FROM binary_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND time_stamp > $2`, userLogin, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prompt, data, note []byte
		var timeStamp time.Time
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

// ForceUpdateCard обновляет информацию о банковской карте.
func (db *DBStorage) ForceUpdateCard(ctx context.Context, userLogin string, prompt []byte,
	number []byte, date []byte, code []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE cards 
		SET prompt = $1, date = $2, code = $3, note = $4, time_stamp = $5
		WHERE user_id = (SELECT user_id FROM users WHERE login = $6)
		AND number = $7`,
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

// ForceUpdateLoginPwd обновляет информацию о паре логин-пароль.
func (db *DBStorage) ForceUpdateLoginPwd(ctx context.Context, userLogin string, prompt []byte,
	login []byte, pwd []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE logins 
		SET pwd = $1, note = $2, time_stamp = $3
		WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
		AND prompt = $5
		AND login = $6`,
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

// ForceUpdateTextRecord обновляет текстовую информацию.
func (db *DBStorage) ForceUpdateTextRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE text_data 
		SET data = $1, note = $2, time_stamp = $3
		WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
		AND prompt = $5`,
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

// ForceUpdateBinaryRecord обновляет бинарные данные.
func (db *DBStorage) ForceUpdateBinaryRecord(ctx context.Context, userLogin string, prompt []byte,
	data []byte, note []byte, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE binary_data 
		SET data = $1, note = $2, time_stamp = $3
		WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
		AND prompt = $5`,
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
