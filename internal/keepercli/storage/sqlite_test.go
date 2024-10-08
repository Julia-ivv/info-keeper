package storage

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	errTest         = errors.New("error")
	testUserLogin   = "ulogin"
	testUserPwd     = "pwd"
	testTime        = "2024-01-02T15:04:05Z"
	testTimeEarlier = "2023-01-02T15:04:05Z"
	testCard        = Card{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Number:    []byte{73, 166, 196, 108, 151, 209, 83, 94, 125, 84, 187, 247, 232, 38, 156, 242, 51, 211, 249},
		Date:      []byte{73, 166, 196, 108, 151, 209, 83, 94, 125, 84, 187, 247, 232, 38, 156, 242, 51, 211, 249},
		Code:      []byte{73, 166, 196, 108, 151, 209, 83, 94, 125, 84, 187, 247, 232, 38, 156, 242, 51, 211, 249},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: testTime,
	}
	testLoginPwd = LoginPwd{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Login:     []byte{20, 251, 144, 177, 213, 73, 247, 129, 45, 118, 66, 54, 135, 0, 6, 123, 43, 112, 216, 226, 196},
		Pwd:       []byte{8, 227, 147, 155, 5, 129, 24, 36, 16, 17, 54, 190, 134, 36, 109, 15, 120, 13, 182},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: testTime,
	}
	testTextRecord = TextRecord{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Data:      []byte{28, 245, 131, 185, 26, 163, 70, 69, 76, 247, 2, 120, 47, 78, 124, 93, 42, 221, 164, 239},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: testTime,
	}
	testBinaryRecord = BinaryRecord{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Data:      []byte{75, 85},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: testTime,
	}
)

func TestCreateTables(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type mockBehavior func()

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			mockBehavior: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS logins").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS cards").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS text_data").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS binary_data").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "create user error",
			mockBehavior: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "create login pwd error",
			mockBehavior: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS logins").WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "create cards error",
			mockBehavior: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS logins").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS cards").WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "create text error",
			mockBehavior: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS logins").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS cards").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS text_data").WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "create binary record error",
			mockBehavior: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS logins").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS cards").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS text_data").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS binary_data").WillReturnError(errTest)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			err := createTables(context.Background(), testDB.dbHandle)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type mockBehavior func()

	tests := []struct {
		name         string
		ctx          context.Context
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "insert error",
			ctx:  context.Background(),
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs([]driver.Value{testUserLogin, testUserPwd, "salt"}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "insert error rows",
			ctx:  context.Background(),
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs([]driver.Value{testUserLogin, testUserPwd, "salt"}...).
					WillReturnResult(sqlmock.NewResult(1, 2))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			err := testDB.RegUser(tt.ctx, testUserLogin, testUserPwd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "wrong hash test",
			ctx:  context.Background(),
			args: args{
				rows:   []string{"hash", "salt"},
				values: []driver.Value{"hash", "salt"},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT hash, salt FROM users").
					WithArgs([]driver.Value{testUserLogin}).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "error select",
			ctx:  context.Background(),
			args: args{
				rows:   []string{"hash"},
				values: []driver.Value{"hash"},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT hash, salt FROM users").
					WithArgs([]driver.Value{testUserLogin}).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.AuthUser(tt.ctx, testUserLogin, testUserPwd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserCardsAfterTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      Card
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      []Card
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testCard,
				rows:   []string{"prompt", "number", "date", "code", "note", "time_stamp"},
				values: []driver.Value{testCard.Prompt, testCard.Number, testCard.Date, testCard.Code, testCard.Note, testCard.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT prompt, number, date, code, note, time_stamp FROM cards").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnRows(rows)
			},
			wantRes: []Card{testCard},
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testCard,
				rows:   []string{"prompt", "number", "date", "code", "note", "time_stamp"},
				values: []driver.Value{testCard.Prompt, testCard.Number, testCard.Date, testCard.Code, testCard.Note, testCard.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT prompt, number, date, code, note, time_stamp FROM cards").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnError(errTest)
			},
			wantRes: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetUserCardsAfterTime(tt.ctx, testUserLogin, testTimeEarlier)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetCard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      Card
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      Card
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testCard,
				rows:   []string{"prompt", "date", "code", "note", "time_stamp"},
				values: []driver.Value{testCard.Prompt, testCard.Date, testCard.Code, testCard.Note, testCard.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT prompt, date, code, note, time_stamp FROM cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Number}...).
					WillReturnRows(rows)
			},
			wantRes: testCard,
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testCard,
				rows:   []string{"prompt", "date", "code", "note", "time_stamp"},
				values: []driver.Value{testCard.Prompt, testCard.Date, testCard.Code, testCard.Note, testCard.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT prompt, date, code, note, time_stamp FROM cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Number}...).
					WillReturnError(errTest)
			},
			wantRes: testCard,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetCard(tt.ctx, testUserLogin, tt.args.c.Number)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetUserLoginsPwdsAfterTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      LoginPwd
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      []LoginPwd
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testLoginPwd,
				rows:   []string{"prompt", "login", "pwd", "note", "time_stamp"},
				values: []driver.Value{testLoginPwd.Prompt, testLoginPwd.Login, testLoginPwd.Pwd, testLoginPwd.Note, testLoginPwd.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT prompt, login, pwd, note, time_stamp FROM logins").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnRows(rows)
			},
			wantRes: []LoginPwd{testLoginPwd},
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testLoginPwd,
				rows:   []string{"prompt", "login", "pwd", "note", "time_stamp"},
				values: []driver.Value{testLoginPwd.Prompt, testLoginPwd.Login, testLoginPwd.Pwd, testLoginPwd.Note, testLoginPwd.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT prompt, login, pwd, note, time_stamp FROM logins").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnError(errTest)
			},
			wantRes: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetUserLoginsPwdsAfterTime(tt.ctx, testUserLogin, testTimeEarlier)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetLoginPwd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      LoginPwd
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      LoginPwd
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testLoginPwd,
				rows:   []string{"pwd", "note", "time_stamp"},
				values: []driver.Value{testLoginPwd.Pwd, testLoginPwd.Note, testLoginPwd.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT pwd, note, time_stamp FROM logins").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Login}...).
					WillReturnRows(rows)
			},
			wantRes: testLoginPwd,
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testLoginPwd,
				rows:   []string{"prompt", "note", "time_stamp"},
				values: []driver.Value{testLoginPwd.Pwd, testLoginPwd.Note, testLoginPwd.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT prompt, date, code, note, time_stamp FROM cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Login}...).
					WillReturnError(errTest)
			},
			wantRes: testLoginPwd,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetLoginPwd(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetUserTextRecordsAfterTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      TextRecord
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      []TextRecord
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testTextRecord,
				rows:   []string{"prompt", "data", "note", "time_stamp"},
				values: []driver.Value{testTextRecord.Prompt, testTextRecord.Data, testTextRecord.Note, testTextRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT prompt, data, note, time_stamp FROM text_data").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnRows(rows)
			},
			wantRes: []TextRecord{testTextRecord},
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testTextRecord,
				rows:   []string{"prompt", "data", "note", "time_stamp"},
				values: []driver.Value{testTextRecord.Prompt, testTextRecord.Data, testTextRecord.Note, testTextRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT prompt, data, note, time_stamp FROM text_data").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnError(errTest)
			},
			wantRes: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetUserTextRecordsAfterTime(tt.ctx, testUserLogin, testTimeEarlier)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetTextRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      TextRecord
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      TextRecord
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testTextRecord,
				rows:   []string{"data", "note", "time_stamp"},
				values: []driver.Value{testTextRecord.Data, testTextRecord.Note, testTextRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT data, note, time_stamp FROM text_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt}...).
					WillReturnRows(rows)
			},
			wantRes: testTextRecord,
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testTextRecord,
				rows:   []string{"data", "note", "time_stamp"},
				values: []driver.Value{testTextRecord.Data, testTextRecord.Note, testTextRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT data, note, time_stamp FROM text_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetTextRecord(tt.ctx, testUserLogin, tt.args.c.Prompt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetUserBinaryRecordAfterTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      BinaryRecord
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      []BinaryRecord
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testBinaryRecord,
				rows:   []string{"prompt", "data", "note", "time_stamp"},
				values: []driver.Value{testBinaryRecord.Prompt, testBinaryRecord.Data, testBinaryRecord.Note, testBinaryRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT prompt, data, note, time_stamp FROM binary_data").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnRows(rows)
			},
			wantRes: []BinaryRecord{testBinaryRecord},
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testBinaryRecord,
				rows:   []string{"prompt", "data", "note", "time_stamp"},
				values: []driver.Value{testBinaryRecord.Prompt, testBinaryRecord.Data, testBinaryRecord.Note, testBinaryRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT data, note, time_stamp FROM binary_data").
					WithArgs([]driver.Value{testUserLogin, testTimeEarlier}...).
					WillReturnError(errTest)
			},
			wantRes: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetUserBinaryRecordsAfterTime(tt.ctx, testUserLogin, testTimeEarlier)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetBinaryRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c      BinaryRecord
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      BinaryRecord
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c:      testBinaryRecord,
				rows:   []string{"data", "note", "time_stamp"},
				values: []driver.Value{testBinaryRecord.Data, testBinaryRecord.Note, testBinaryRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT data, note, time_stamp FROM binary_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt}...).
					WillReturnRows(rows)
			},
			wantRes: testBinaryRecord,
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c:      testBinaryRecord,
				rows:   []string{"data", "note", "time_stamp"},
				values: []driver.Value{testBinaryRecord.Data, testBinaryRecord.Note, testBinaryRecord.TimeStamp},
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT data, note, time_stamp FROM binary_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt}...).
					WillReturnError(errTest)
			},
			wantRes: testBinaryRecord,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			c, err := testDB.GetBinaryRecord(tt.ctx, testUserLogin, tt.args.c.Prompt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, c)
			}
		})
	}
}

func TestGetLastSyncTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		rows   []string
		values []driver.Value
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantRes      string
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				rows:   []string{"last_sync"},
				values: []driver.Value{testTime},
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows(a.rows).AddRow(a.values...)
				mock.ExpectQuery("SELECT last_sync FROM users").
					WithArgs([]driver.Value{testUserLogin}...).
					WillReturnRows(rows)
			},
			wantRes: testTime,
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{},
			mockBehavior: func(a args) {
				mock.ExpectQuery("SELECT last_sync FROM users").
					WithArgs([]driver.Value{testUserLogin}...).
					WillReturnError(errTest)
			},
			wantRes: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			l, err := testDB.GetLastSyncTime(tt.ctx, testUserLogin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, l)
			}
		})
	}
}

func TestUpdateLastSyncTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type mockBehavior func()

	tests := []struct {
		name         string
		ctx          context.Context
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			mockBehavior: func() {
				mock.ExpectExec("UPDATE users").
					WithArgs([]driver.Value{testTime, testUserLogin}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error",
			ctx:  context.Background(),
			mockBehavior: func() {
				mock.ExpectQuery("UPDATE users").
					WithArgs([]driver.Value{testTime, testUserLogin}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			mockBehavior: func() {
				mock.ExpectExec("UPDATE users").
					WithArgs([]driver.Value{testTime, testUserLogin}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			err := testDB.UpdateLastSyncTime(tt.ctx, testUserLogin, testTime)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddCard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c Card
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testCard,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date,
						a.c.Code, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testCard,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date,
						a.c.Code, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date,
						a.c.Code, a.c.Note, a.c.TimeStamp}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.AddCard(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Number,
				tt.args.c.Date, tt.args.c.Code, tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddLoginPwd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c LoginPwd
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testLoginPwd,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testLoginPwd,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c: testLoginPwd,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, a.c.TimeStamp}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.AddLoginPwd(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Login,
				tt.args.c.Pwd, tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddTextRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c TextRecord
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testTextRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Data, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testTextRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Data, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c: testTextRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Data, a.c.Note, a.c.TimeStamp}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.AddTextRecord(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Data, tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddBinaryRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c BinaryRecord
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO binary_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Data, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO binary_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Data, a.c.Note, a.c.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  context.Background(),
			args: args{
				c: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("INSERT INTO binary_data").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Data, a.c.Note, a.c.TimeStamp}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.AddBinaryRecord(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Data,
				tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddSyncData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c Card
		l LoginPwd
		t TextRecord
		b BinaryRecord
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.t.Prompt, a.t.Data, a.t.Note, a.t.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO binary_data").
					WithArgs([]driver.Value{testUserLogin, a.b.Prompt, a.b.Data, a.b.Note, a.b.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "insert card error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnError(errTest)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert card row error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert login error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnError(errTest)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert login row error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert text error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.t.Prompt, a.t.Data, a.t.Note, a.t.TimeStamp}...).
					WillReturnError(errTest)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert text row error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.t.Prompt, a.t.Data, a.t.Note, a.t.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert bytes error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.t.Prompt, a.t.Data, a.t.Note, a.t.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO binary_data").
					WithArgs([]driver.Value{testUserLogin, a.b.Prompt, a.b.Data, a.b.Note, a.b.TimeStamp}...).
					WillReturnError(errTest)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "insert bytes row error",
			ctx:  context.Background(),
			args: args{
				c: testCard,
				l: testLoginPwd,
				t: testTextRecord,
				b: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO cards").
					WithArgs([]driver.Value{testUserLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code,
						a.c.Note, a.c.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO logins").
					WithArgs([]driver.Value{testUserLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note,
						a.l.TimeStamp}...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO text_data").
					WithArgs([]driver.Value{testUserLogin, a.t.Prompt, a.t.Data, a.t.Note, a.t.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO binary_data").
					WithArgs([]driver.Value{testUserLogin, a.b.Prompt, a.b.Data, a.b.Note, a.b.TimeStamp}...).
					WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.AddSyncData(tt.ctx, testUserLogin, []Card{testCard}, []LoginPwd{testLoginPwd},
				[]TextRecord{testTextRecord}, []BinaryRecord{testBinaryRecord})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateCard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c Card
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testCard,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE cards").
					WithArgs([]driver.Value{a.c.Prompt, a.c.Date, a.c.Code, a.c.Note,
						a.c.TimeStamp, testUserLogin, a.c.Number}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error update",
			ctx:  context.Background(),
			args: args{
				c: testCard,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE cards").
					WithArgs([]driver.Value{a.c.Prompt, a.c.Date, a.c.Code, a.c.Note,
						a.c.TimeStamp, testUserLogin, a.c.Number}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testCard,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE cards").
					WithArgs([]driver.Value{a.c.Prompt, a.c.Date, a.c.Code, a.c.Note,
						a.c.TimeStamp, testUserLogin, a.c.Number}...).
					WillReturnResult(sqlmock.NewResult(1, 2))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.UpdateCard(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Number, tt.args.c.Date,
				tt.args.c.Code, tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c LoginPwd
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testLoginPwd,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE logins").
					WithArgs([]driver.Value{a.c.Pwd, a.c.Note, a.c.TimeStamp, testUserLogin,
						a.c.Prompt, a.c.Login}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error update",
			ctx:  context.Background(),
			args: args{
				c: testLoginPwd,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE logins").
					WithArgs([]driver.Value{a.c.Pwd, a.c.Note, a.c.TimeStamp, testUserLogin,
						a.c.Prompt, a.c.Login}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testLoginPwd,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE logins").
					WithArgs([]driver.Value{a.c.Pwd, a.c.Note, a.c.TimeStamp, testUserLogin,
						a.c.Prompt, a.c.Login}...).
					WillReturnResult(sqlmock.NewResult(1, 2))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.UpdateLoginPwd(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Login,
				tt.args.c.Pwd, tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateTextRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c TextRecord
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testTextRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE text_data").
					WithArgs([]driver.Value{a.c.Data, a.c.Note, a.c.TimeStamp, testUserLogin, a.c.Prompt}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error update",
			ctx:  context.Background(),
			args: args{
				c: testTextRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE text_data").
					WithArgs([]driver.Value{a.c.Data, a.c.Note, a.c.TimeStamp, testUserLogin, a.c.Prompt}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testTextRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE text_data").
					WithArgs([]driver.Value{a.c.Data, a.c.Note, a.c.TimeStamp, testUserLogin, a.c.Prompt}...).
					WillReturnResult(sqlmock.NewResult(1, 2))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.UpdateTextRecord(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Data,
				tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateBinaryRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock: %s", err)
	}
	defer db.Close()

	testDB := SQLiteStorage{dbHandle: db}

	type args struct {
		c BinaryRecord
	}
	type mockBehavior func(a args)

	tests := []struct {
		name         string
		ctx          context.Context
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			args: args{
				c: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE binary_data").
					WithArgs([]driver.Value{a.c.Data, a.c.Note, a.c.TimeStamp, testUserLogin, a.c.Prompt}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error update",
			ctx:  context.Background(),
			args: args{
				c: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE binary_data").
					WithArgs([]driver.Value{a.c.Data, a.c.Note, a.c.TimeStamp, testUserLogin, a.c.Prompt}...).
					WillReturnError(errTest)
			},
			wantErr: true,
		},
		{
			name: "error rows",
			ctx:  context.Background(),
			args: args{
				c: testBinaryRecord,
			},
			mockBehavior: func(a args) {
				mock.ExpectExec("UPDATE binary_data").
					WithArgs([]driver.Value{a.c.Data, a.c.Note, a.c.TimeStamp, testUserLogin, a.c.Prompt}...).
					WillReturnResult(sqlmock.NewResult(1, 2))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := testDB.UpdateBinaryRecord(tt.ctx, testUserLogin, tt.args.c.Prompt, tt.args.c.Data,
				tt.args.c.Note, tt.args.c.TimeStamp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
