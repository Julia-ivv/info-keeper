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
	"github.com/Julia-ivv/info-keeper.git/internal/cryptor"
	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
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
			user_id serial UNIQUE, 
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
			prompt bytea NOT NULL,
			login bytea NOT NULL,
			pwd bytea NOT NULL,
			note bytea,
			time_stamp timestamptz (0) NOT NULL,
			PRIMARY KEY(user_id, login, prompt)
		)`)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return &DBStorage{dbHandle: db}, nil
}

func (db *DBStorage) Close() error {
	return db.dbHandle.Close()
}

func (db *DBStorage) RegUser(ctx context.Context, login string, pwd string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	salt, err := GenerateRandomString(LengthSalt)
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

func (db *DBStorage) AddCard(ctx context.Context, userLogin string, prompt string,
	number string, date string, code string, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || number == "" || date == "" || code == "" {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNumber, err := cryptor.EncryptsString(number)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encDate, err := cryptor.EncryptsString(date)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encCode, err := cryptor.EncryptsString(code)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO cards (user_id , prompt, number, date, code, note, time_stamp) 
			VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5, $6, $7)`,
		userLogin, encPrompt, encNumber, encDate, encCode, encNote, timeStamp)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM cards 
				WHERE number = $1 AND 
				user_id = (SELECT user_id FROM users WHERE login = $2)`, encNumber, userLogin)
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
				encPrompt, encDate, encCode, encNote, timeStamp, userLogin, encNumber)
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

func (db *DBStorage) AddLoginPwd(ctx context.Context, userLogin string, prompt string,
	login string, pwd string, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || login == "" || pwd == "" {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encLogin, err := cryptor.EncryptsString(login)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encPwd, err := cryptor.EncryptsString(pwd)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO logins (user_id , prompt, login, pwd, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5, $6)`,
		userLogin, encPrompt, encLogin, encPwd, encNote, timeStamp)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM logins 
				WHERE prompt = $1 
				AND login = $2 
				AND user_id = (SELECT user_id FROM users WHERE login = $3)`,
				encPrompt, encLogin, userLogin)
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
				encPwd, encNote, timeStamp, userLogin, encPrompt, encLogin)
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

func (db *DBStorage) AddTextRecord(ctx context.Context, userLogin string, prompt string,
	data string, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || data == "" {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encData, err := cryptor.EncryptsString(data)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO text_data (user_id , prompt, data, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5)`,
		userLogin, encPrompt, encData, encNote, timeStamp)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM text_data 
				WHERE prompt = $1 AND 
				user_id = (SELECT user_id FROM users WHERE login = $2)`, encPrompt, userLogin)
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
				encData, encNote, timeStamp, userLogin, encPrompt)
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

func (db *DBStorage) AddBinaryRecord(ctx context.Context, userLogin string, prompt string,
	data []byte, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || data == nil {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encData, err := cryptor.EncryptsByte(data)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`INSERT INTO binary_data (user_id , prompt, data, note, time_stamp) 
		VALUES ((SELECT user_id FROM users WHERE login = $1), $2, $3, $4, $5)`,
		userLogin, encPrompt, encData, encNote, timeStamp)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.dbHandle.QueryRowContext(ctx,
				`SELECT time_stamp FROM binary_data 
				WHERE prompt = $1 AND 
				user_id = (SELECT user_id FROM users WHERE login = $2)`, encPrompt, userLogin)
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
				encData, encNote, timeStamp, userLogin, encPrompt)
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

type Card struct {
	Prompt    string
	Number    string
	Date      string
	Code      string
	Note      string
	TimeStamp time.Time
}

func (db *DBStorage) GetCard(ctx context.Context, userLogin string, number string) (card Card, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	enNumber, err := cryptor.EncryptsString(number)
	if err != nil {
		return Card{}, NewStorError(EncryptionError, err)
	}
	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT prompt, date, code, note, time_stamp
		FROM cards
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND number = $2`, userLogin, enNumber)

	var prompt, date, code, note []byte
	var timeStamp time.Time
	err = row.Scan(&prompt, &date, &code, &note, &timeStamp)
	if err != nil {
		return Card{}, err
	}

	decPrompt, err := cryptor.Decrypts(prompt)
	if err != nil {
		return Card{}, NewStorError(DecryptionError, err)
	}
	decDate, err := cryptor.Decrypts(date)
	if err != nil {
		return Card{}, NewStorError(DecryptionError, err)
	}
	decCode, err := cryptor.Decrypts(code)
	if err != nil {
		return Card{}, NewStorError(DecryptionError, err)
	}
	decNote, err := cryptor.Decrypts(note)
	if err != nil {
		return Card{}, NewStorError(DecryptionError, err)
	}

	return Card{
		Prompt:    decPrompt,
		Number:    number,
		Date:      decDate,
		Code:      decCode,
		Note:      decNote,
		TimeStamp: timeStamp,
	}, nil
}

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
		decPrompt, err := cryptor.Decrypts(prompt)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decNumber, err := cryptor.Decrypts(number)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decDate, err := cryptor.Decrypts(date)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decCode, err := cryptor.Decrypts(code)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decNote, err := cryptor.Decrypts(note)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		cards = append(cards, Card{
			Prompt:    decPrompt,
			Number:    decNumber,
			Date:      decDate,
			Code:      decCode,
			Note:      decNote,
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
	Prompt    string
	Login     string
	Pwd       string
	Note      string
	TimeStamp time.Time
}

func (db *DBStorage) GetLoginPwd(ctx context.Context, userLogin string, prompt string, login string) (loginPwd LoginPwd, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	enPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return LoginPwd{}, NewStorError(EncryptionError, err)
	}
	enLogin, err := cryptor.EncryptsString(login)
	if err != nil {
		return LoginPwd{}, NewStorError(EncryptionError, err)
	}
	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT pwd, note, time_stamp
		FROM logins
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND prompt = $2 AND login = $3`, userLogin, enPrompt, enLogin)

	var pwd, note []byte
	var timeStamp time.Time
	err = row.Scan(&pwd, &note, &timeStamp)
	if err != nil {
		return LoginPwd{}, err
	}

	decPwd, err := cryptor.Decrypts(pwd)
	if err != nil {
		return LoginPwd{}, NewStorError(DecryptionError, err)
	}
	decNote, err := cryptor.Decrypts(note)
	if err != nil {
		return LoginPwd{}, NewStorError(DecryptionError, err)
	}

	return LoginPwd{
		Prompt:    prompt,
		Login:     login,
		Pwd:       decPwd,
		Note:      decNote,
		TimeStamp: timeStamp,
	}, nil
}

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
		decPrompt, err := cryptor.Decrypts(prompt)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decLogin, err := cryptor.Decrypts(login)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decPwd, err := cryptor.Decrypts(pwd)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decNote, err := cryptor.Decrypts(note)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		loginsPwds = append(loginsPwds, LoginPwd{
			Prompt:    decPrompt,
			Login:     decLogin,
			Pwd:       decPwd,
			Note:      decNote,
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
	Prompt    string
	Data      string
	Note      string
	TimeStamp time.Time
}

func (db *DBStorage) GetTextRecord(ctx context.Context, userLogin string, prompt string) (record TextRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	enPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return TextRecord{}, NewStorError(EncryptionError, err)
	}
	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT data, note, time_stamp
		FROM text_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND prompt = $2`, userLogin, enPrompt)

	var data, note []byte
	var timeStamp time.Time
	err = row.Scan(&data, &note, &timeStamp)
	if err != nil {
		return TextRecord{}, err
	}

	decData, err := cryptor.Decrypts(data)
	if err != nil {
		return TextRecord{}, NewStorError(DecryptionError, err)
	}
	decNote, err := cryptor.Decrypts(note)
	if err != nil {
		return TextRecord{}, NewStorError(DecryptionError, err)
	}

	return TextRecord{
		Prompt:    prompt,
		Data:      decData,
		Note:      decNote,
		TimeStamp: timeStamp,
	}, nil
}

// func (db *DBStorage) GetUserTextPrompts(ctx context.Context, userLogin string) (prompts []string, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()
// 	//
// 	rows, err := db.dbHandle.QueryContext(ctx,
// 		`SELECT prompt
// 		FROM text_data
// 		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)`, userLogin)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var p string
// 		err = rows.Scan(&p)
// 		if err != nil {
// 			return nil, err
// 		}
// 		prompts = append(prompts, p)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return prompts, nil
// }

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
		decPrompt, err := cryptor.Decrypts(prompt)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decData, err := cryptor.Decrypts(data)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decNote, err := cryptor.Decrypts(note)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		records = append(records, TextRecord{
			Prompt:    decPrompt,
			Data:      decData,
			Note:      decNote,
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
	Prompt    string
	Data      []byte
	Note      string
	TimeStamp time.Time
}

func (db *DBStorage) GetBinaryRecord(ctx context.Context, userLogin string, prompt string) (record BinaryRecord, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	enPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return BinaryRecord{}, NewStorError(EncryptionError, err)
	}
	row := db.dbHandle.QueryRowContext(ctx,
		`SELECT data, note, time_stamp
		FROM binary_data
		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)
		AND prompt = $2`, userLogin, enPrompt)

	var data, note []byte
	var timeStamp time.Time
	err = row.Scan(&data, &note, &timeStamp)
	if err != nil {
		return BinaryRecord{}, err
	}

	decData, err := cryptor.DecryptsByte(data)
	if err != nil {
		return BinaryRecord{}, NewStorError(DecryptionError, err)
	}
	decNote, err := cryptor.Decrypts(note)
	if err != nil {
		return BinaryRecord{}, NewStorError(DecryptionError, err)
	}

	logger.ZapSugar.Info("record data ", record.Data)

	return BinaryRecord{
		Prompt:    prompt,
		Data:      decData,
		Note:      decNote,
		TimeStamp: timeStamp,
	}, nil
}

// func (db *DBStorage) GetUserBinaryPrompts(ctx context.Context, userLogin string) (prompts []string, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()
// 	//
// 	rows, err := db.dbHandle.QueryContext(ctx,
// 		`SELECT prompt
// 		FROM binary_data
// 		WHERE user_id = (SELECT user_id FROM users WHERE login = $1)`, userLogin)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var p string
// 		err = rows.Scan(&p)
// 		if err != nil {
// 			return nil, err
// 		}
// 		prompts = append(prompts, p)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return prompts, nil
// }

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
		decPrompt, err := cryptor.Decrypts(prompt)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decData, err := cryptor.DecryptsByte(data)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		decNote, err := cryptor.Decrypts(note)
		if err != nil {
			return nil, NewStorError(DecryptionError, err)
		}
		records = append(records, BinaryRecord{
			Prompt:    decPrompt,
			Data:      decData,
			Note:      decNote,
			TimeStamp: timeStamp,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (db *DBStorage) ForceUpdateCard(ctx context.Context, userLogin string, prompt string,
	number string, date string, code string, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || number == "" || date == "" || code == "" {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNumber, err := cryptor.EncryptsString(number)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encDate, err := cryptor.EncryptsString(date)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encCode, err := cryptor.EncryptsString(code)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE cards 
		SET prompt = $1, date = $2, code = $3, note = $4, time_stamp = $5
		WHERE user_id = (SELECT user_id FROM users WHERE login = $6)
		AND number = $7`,
		encPrompt, encDate, encCode, encNote, timeStamp, userLogin, encNumber)
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

func (db *DBStorage) ForceUpdateLoginPwd(ctx context.Context, userLogin string, prompt string,
	login string, pwd string, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || login == "" || pwd == "" {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encLogin, err := cryptor.EncryptsString(login)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encPwd, err := cryptor.EncryptsString(pwd)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE logins 
		SET pwd = $1, note = $2, time_stamp = $3
		WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
		AND prompt = $5
		AND login = $6`,
		encPwd, encNote, timeStamp, userLogin, encPrompt, encLogin)
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

func (db *DBStorage) ForceUpdateTextRecord(ctx context.Context, userLogin string, prompt string,
	data string, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || data == "" {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encData, err := cryptor.EncryptsString(data)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE text_data 
		SET data = $1, note = $2, time_stamp = $3
		WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
		AND prompt = $5`,
		encData, encNote, timeStamp, userLogin, encPrompt)
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

func (db *DBStorage) ForceUpdateBinaryRecord(ctx context.Context, userLogin string, prompt string,
	data []byte, note string, timeStamp time.Time) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if prompt == "" || data == nil {
		return NewStorError(EmptyValues, errors.New("empty required fields"))
	}

	encPrompt, err := cryptor.EncryptsString(prompt)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encData, err := cryptor.EncryptsByte(data)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}
	encNote, err := cryptor.EncryptsString(note)
	if err != nil {
		return NewStorError(EncryptionError, err)
	}

	result, err := db.dbHandle.ExecContext(ctx,
		`UPDATE binary_data 
		SET data = $1, note = $2, time_stamp = $3
		WHERE user_id = (SELECT user_id FROM users WHERE login = $4)
		AND prompt = $5`,
		encData, encNote, timeStamp, userLogin, encPrompt)
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
