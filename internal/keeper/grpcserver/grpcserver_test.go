package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keeper/config"
	"github.com/Julia-ivv/info-keeper.git/internal/keeper/mocks"
	"github.com/Julia-ivv/info-keeper.git/internal/keeper/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

const (
	testTime = "2024-01-02T15:04:05Z"
)

var (
	testLoginPwd = storage.LoginPwd{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Login:     []byte{20, 251, 144, 177, 213, 73, 247, 129, 45, 118, 66, 54, 135, 0, 6, 123, 43, 112, 216, 226, 196},
		Pwd:       []byte{8, 227, 147, 155, 5, 129, 24, 36, 16, 17, 54, 190, 134, 36, 109, 15, 120, 13, 182},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: time.Time{},
	}
	testCard = storage.Card{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Number:    []byte{73, 166, 196, 108, 151, 209, 83, 94, 125, 84, 187, 247, 232, 38, 156, 242, 51, 211, 249},
		Date:      []byte{73, 166, 196, 108, 151, 209, 83, 94, 125, 84, 187, 247, 232, 38, 156, 242, 51, 211, 249},
		Code:      []byte{73, 166, 196, 108, 151, 209, 83, 94, 125, 84, 187, 247, 232, 38, 156, 242, 51, 211, 249},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: time.Time{},
	}
	testTextRecord = storage.TextRecord{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Data:      []byte{28, 245, 131, 185, 26, 163, 70, 69, 76, 247, 2, 120, 47, 78, 124, 93, 42, 221, 164, 239},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: time.Time{},
	}
	testBinaryRecord = storage.BinaryRecord{
		Prompt:    []byte{8, 230, 152, 2, 249, 163, 40, 83, 43, 16, 152, 201, 204, 108, 25, 36, 123, 91, 33},
		Data:      []byte{75, 85},
		Note:      []byte{22, 251, 131, 189, 215, 11, 255, 110, 100, 134, 201, 112, 212, 94, 71, 85, 234, 240, 187, 114},
		TimeStamp: time.Time{},
	}
	testCardPb = &pb.UserCard{
		Prompt:    testCard.Prompt,
		Number:    testCard.Number,
		Date:      testCard.Date,
		Code:      testCard.Code,
		Note:      testCard.Note,
		TimeStamp: time.Time{}.Format(time.RFC3339),
	}
	testLoginPwdPb = &pb.UserLoginPwd{
		Prompt:    testLoginPwd.Prompt,
		Login:     testLoginPwd.Login,
		Pwd:       testLoginPwd.Pwd,
		Note:      testLoginPwd.Note,
		TimeStamp: time.Time{}.Format(time.RFC3339),
	}
	testTextPb = &pb.UserTextRecord{
		Prompt:    testTextRecord.Prompt,
		Data:      testTextRecord.Data,
		Note:      testTextRecord.Note,
		TimeStamp: time.Time{}.Format(time.RFC3339),
	}
	testBinaryPb = &pb.UserBinaryRecord{
		Prompt:    testBinaryRecord.Prompt,
		Data:      testBinaryRecord.Data,
		Note:      testBinaryRecord.Note,
		TimeStamp: time.Time{}.Format(time.RFC3339),
	}
	testUserLogin = "ulogin"
	testUserPwd   = "ulogin"
	testCfg       = config.Flags{SecretKey: "rtyhg"}
)

func TestAddUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
		pwd   string
	}

	tests := []struct {
		name      string
		prepare   func(m *mocks.MockRepositorier, a args)
		args      args
		expectRes *pb.AddUserResponse
		wantErr   bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				gomock.InOrder(
					m.EXPECT().RegUser(a.ctx, a.login, a.pwd).Return(nil),
					m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(nil),
				)
			},
			args: args{
				ctx:   context.Background(),
				login: "user1",
				pwd:   "pwd1",
			},
			expectRes: &pb.AddUserResponse{
				Token: "some-token",
			},
			wantErr: false,
		},
		{
			name: "empty data test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				gomock.InOrder(
					m.EXPECT().RegUser(a.ctx, a.login, a.pwd).Return(nil).AnyTimes(),
					m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:   context.Background(),
				login: "",
				pwd:   "",
			},
			expectRes: nil,
			wantErr:   true,
		},
		{
			name: "error registration test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				gomock.InOrder(
					m.EXPECT().RegUser(a.ctx, a.login, a.pwd).Return(errors.New("")),
					m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(errors.New("")).AnyTimes(),
				)
			},
			args: args{
				ctx:   context.Background(),
				login: "user",
				pwd:   "pwd",
			},
			expectRes: nil,
			wantErr:   true,
		},
		{
			name: "error auth test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				gomock.InOrder(
					m.EXPECT().RegUser(a.ctx, a.login, a.pwd).Return(nil),
					m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(errors.New("")).AnyTimes(),
				)
			},
			args: args{
				ctx:   context.Background(),
				login: "user",
				pwd:   "pwd",
			},
			expectRes: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, config.Flags{})
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.AddUser(tt.args.ctx, &pb.AddUserRequest{
				Login: tt.args.login,
				Pwd:   tt.args.pwd,
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			}
		})
	}
}

func TestAuthUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
		pwd   string
	}

	tests := []struct {
		name      string
		prepare   func(m *mocks.MockRepositorier, a args)
		args      args
		expectRes *pb.AddUserResponse
		wantErr   bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(nil)
			},
			args: args{
				ctx:   context.Background(),
				login: "user1",
				pwd:   "pwd1",
			},
			expectRes: &pb.AddUserResponse{
				Token: "some-token",
			},
			wantErr: false,
		},
		{
			name: "empty data test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(nil).AnyTimes()
			},
			args: args{
				ctx:   context.Background(),
				login: "",
				pwd:   "",
			},
			expectRes: nil,
			wantErr:   true,
		},
		{
			name: "error auth test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().AuthUser(a.ctx, a.login, a.pwd).Return(errors.New("error auth"))
			},
			args: args{
				ctx:   context.Background(),
				login: "user",
				pwd:   "pwd",
			},
			expectRes: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, config.Flags{})
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.AuthUser(tt.args.ctx, &pb.AuthUserRequest{
				Login: tt.args.login,
				Pwd:   tt.args.pwd,
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			}
		})
	}
}

func TestAddCard(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.Card
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.Card{
					Prompt:    nil,
					Number:    nil,
					Date:      nil,
					Code:      nil,
					Note:      nil,
					TimeStamp: time.Time{},
				},
				timeSt: testTime,
			},
			wantErr: true,
		},
		{
			name: "exists newer test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).
					Return(storage.NewStorError(storage.ExistsDataNewerVersion, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time stamp test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    "2006-01T15:04:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.AddCard(tt.args.ctx, &pb.AddCardRequest{
				Card: &pb.UserCard{
					Prompt:    tt.args.c.Prompt,
					Number:    tt.args.c.Number,
					Date:      tt.args.c.Date,
					Code:      tt.args.c.Code,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddLogin(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.LoginPwd
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.LoginPwd{
					Prompt:    nil,
					Login:     nil,
					Pwd:       nil,
					Note:      nil,
					TimeStamp: time.Time{},
				},
				timeSt: testTime,
			},
			wantErr: true,
		},
		{
			name: "exists newer test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(storage.NewStorError(storage.ExistsDataNewerVersion, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time stamp test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    "2006-01-02T154:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.AddLogin(tt.args.ctx, &pb.AddLoginRequest{
				LoginPwd: &pb.UserLoginPwd{
					Prompt:    tt.args.c.Prompt,
					Login:     tt.args.c.Login,
					Pwd:       tt.args.c.Pwd,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddTextData(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.TextRecord
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.TextRecord{
					Prompt:    nil,
					Data:      nil,
					Note:      nil,
					TimeStamp: time.Time{},
				},
				timeSt: testTime,
			},
			wantErr: true,
		},
		{
			name: "exists newer test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.ExistsDataNewerVersion, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time parse test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    "2006-01T15:04:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.AddTextData(tt.args.ctx, &pb.AddTextDataRequest{
				TextRecord: &pb.UserTextRecord{
					Prompt:    tt.args.c.Prompt,
					Data:      tt.args.c.Data,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddBinaryData(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.BinaryRecord
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "exists newer test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.ExistsDataNewerVersion, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time parse test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    "2006-01T15:04:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.AddBinaryData(tt.args.ctx, &pb.AddBinaryDataRequest{
				BinaryRecord: &pb.UserBinaryRecord{
					Prompt:    tt.args.c.Prompt,
					Data:      tt.args.c.Data,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSyncUserData(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	testTimePrs, err := time.Parse(time.RFC3339, testTime)
	if err != nil {
		fmt.Println("time parse error")
	}

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.Card
		l         storage.LoginPwd
		t         storage.TextRecord
		b         storage.BinaryRecord
		lastSync  string
	}

	tests := []struct {
		name       string
		prepare    func(m *mocks.MockRepositorier, a args)
		args       args
		inCards    []*pb.UserCard
		inLogins   []*pb.UserLoginPwd
		inTexts    []*pb.UserTextRecord
		inBinaryes []*pb.UserBinaryRecord
		wantRes    *pb.SyncUserDataResponse
		wantErr    bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes: &pb.SyncUserDataResponse{
				SyncErrors:       []*pb.SyncUserDataResponse_SyncErrorInfo{},
				NewLogins:        []*pb.UserLoginPwd{testLoginPwdPb},
				NewCards:         []*pb.UserCard{testCardPb},
				NewTextRecords:   []*pb.UserTextRecord{testTextPb},
				NewBinaryRecords: []*pb.UserBinaryRecord{testBinaryPb},
			},
			wantErr: false,
		},
		{
			name: "ok test with clients data",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c, a.c}, nil).AnyTimes(),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l, a.l}, nil).AnyTimes(),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t, a.t}, nil).AnyTimes(),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b, a.b}, nil).AnyTimes(),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  "0001-01-01T00:00:00Z",
			},
			inCards:    []*pb.UserCard{testCardPb},
			inLogins:   []*pb.UserLoginPwd{testLoginPwdPb},
			inTexts:    []*pb.UserTextRecord{testTextPb},
			inBinaryes: []*pb.UserBinaryRecord{testBinaryPb},
			wantRes: &pb.SyncUserDataResponse{
				SyncErrors:       []*pb.SyncUserDataResponse_SyncErrorInfo{},
				NewLogins:        []*pb.UserLoginPwd{},
				NewCards:         []*pb.UserCard{},
				NewTextRecords:   []*pb.UserTextRecord{},
				NewBinaryRecords: []*pb.UserBinaryRecord{},
			},
			wantErr: false,
		},
		{
			name: "empty user test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil).AnyTimes(),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil).AnyTimes(),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil).AnyTimes(),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil).AnyTimes(),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes: &pb.SyncUserDataResponse{
				SyncErrors:       nil,
				NewLogins:        []*pb.UserLoginPwd{},
				NewCards:         []*pb.UserCard{},
				NewTextRecords:   []*pb.UserTextRecord{},
				NewBinaryRecords: []*pb.UserBinaryRecord{},
			},
			wantErr: true,
		},
		{
			name: "error parse time test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, testTimePrs).
						Return([]storage.Card{a.c}, nil).AnyTimes(),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, testTimePrs).
						Return([]storage.LoginPwd{a.l}, nil).AnyTimes(),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, testTimePrs).
						Return([]storage.TextRecord{a.t}, nil).AnyTimes(),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, testTimePrs).
						Return([]storage.BinaryRecord{a.b}, nil).AnyTimes(),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  "2024-01T15:04:05Z",
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes: &pb.SyncUserDataResponse{
				SyncErrors:       nil,
				NewLogins:        []*pb.UserLoginPwd{},
				NewCards:         []*pb.UserCard{},
				NewTextRecords:   []*pb.UserTextRecord{},
				NewBinaryRecords: []*pb.UserBinaryRecord{},
			},
			wantErr: true,
		},
		{
			name: "get user text after time error",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil).AnyTimes(),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil).AnyTimes(),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return(nil, errors.New("error")),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil).AnyTimes(),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes:    &pb.SyncUserDataResponse{},
			wantErr:    true,
		},
		{
			name: "get user logins after time error",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil).AnyTimes(),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return(nil, errors.New("error")),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil).AnyTimes(),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil).AnyTimes(),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes:    &pb.SyncUserDataResponse{},
			wantErr:    true,
		},
		{
			name: "get user cards after time error",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return(nil, errors.New("error")),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil).AnyTimes(),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil).AnyTimes(),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil).AnyTimes(),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes:    &pb.SyncUserDataResponse{},
			wantErr:    true,
		},
		{
			name: "get user binary after time error",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return(nil, errors.New("error")),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards:    []*pb.UserCard{},
			inLogins:   []*pb.UserLoginPwd{},
			inTexts:    []*pb.UserTextRecord{},
			inBinaryes: []*pb.UserBinaryRecord{},
			wantRes:    &pb.SyncUserDataResponse{},
			wantErr:    true,
		},
		{
			name: "add cards time parse error",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(nil).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return(nil).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  testTime,
			},
			inCards: []*pb.UserCard{{
				Prompt:    testCard.Prompt,
				Number:    testCard.Number,
				Date:      testCard.Date,
				Code:      testCard.Code,
				Note:      testCard.Note,
				TimeStamp: "1",
			}},
			inLogins: []*pb.UserLoginPwd{{
				Prompt:    testLoginPwd.Prompt,
				Login:     testLoginPwd.Login,
				Pwd:       testLoginPwd.Pwd,
				Note:      testLoginPwd.Note,
				TimeStamp: "1",
			}},
			inTexts: []*pb.UserTextRecord{{
				Prompt:    testTextRecord.Prompt,
				Data:      testTextRecord.Data,
				Note:      testTextRecord.Note,
				TimeStamp: "1",
			}},
			inBinaryes: []*pb.UserBinaryRecord{{
				Prompt:    testBinaryRecord.Prompt,
				Data:      testBinaryRecord.Data,
				Note:      testBinaryPb.Note,
				TimeStamp: "1",
			}},
			wantRes: &pb.SyncUserDataResponse{
				SyncErrors: []*pb.SyncUserDataResponse_SyncErrorInfo{
					{
						Text:  "error for card number ",
						Value: testCard.Number,
						Err:   "parsing time \"1\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"1\" as \"2006\"",
					},
					{
						Text:  "error for pair login/password with prompt ",
						Value: testLoginPwd.Prompt,
						Err:   "parsing time \"1\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"1\" as \"2006\"",
					},
					{
						Text:  "error for text data with prompt ",
						Value: testTextRecord.Prompt,
						Err:   "parsing time \"1\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"1\" as \"2006\"",
					},
					{
						Text:  "error for binary data with prompt ",
						Value: testBinaryPb.Prompt,
						Err:   "parsing time \"1\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"1\" as \"2006\"",
					},
				},
				NewLogins:        []*pb.UserLoginPwd{testLoginPwdPb},
				NewCards:         []*pb.UserCard{testCardPb},
				NewTextRecords:   []*pb.UserTextRecord{testTextPb},
				NewBinaryRecords: []*pb.UserBinaryRecord{testBinaryPb},
			},
			wantErr: false,
		},
		{
			name: "add records error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.lastSync)
				require.NoError(t, err)
				gomock.InOrder(
					m.EXPECT().GetUserCardsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.Card{a.c}, nil),
					m.EXPECT().GetUserLoginsPwdsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.LoginPwd{a.l}, nil),
					m.EXPECT().GetUserTextRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.TextRecord{a.t}, nil),
					m.EXPECT().GetUserBinaryRecordsAfterTime(a.ctx, a.userLogin, tp).
						Return([]storage.BinaryRecord{a.b}, nil),
					m.EXPECT().AddCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number, a.c.Date, a.c.Code, a.c.Note, tp).
						Return(errors.New("add card error")).AnyTimes(),
					m.EXPECT().AddLoginPwd(a.ctx, a.userLogin, a.l.Prompt, a.l.Login, a.l.Pwd, a.l.Note, tp).
						Return(errors.New("add login error")).AnyTimes(),
					m.EXPECT().AddTextRecord(a.ctx, a.userLogin, a.t.Prompt, a.t.Data, a.t.Note, tp).
						Return(errors.New("add text error")).AnyTimes(),
					m.EXPECT().AddBinaryRecord(a.ctx, a.userLogin, a.b.Prompt, a.b.Data, a.b.Note, tp).
						Return((errors.New("add bytes error"))).AnyTimes(),
				)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				l:         testLoginPwd,
				t:         testTextRecord,
				b:         testBinaryRecord,
				lastSync:  time.Time{}.Format(time.RFC3339),
			},
			inCards:    []*pb.UserCard{testCardPb},
			inLogins:   []*pb.UserLoginPwd{testLoginPwdPb},
			inTexts:    []*pb.UserTextRecord{testTextPb},
			inBinaryes: []*pb.UserBinaryRecord{testBinaryPb},
			wantRes: &pb.SyncUserDataResponse{
				SyncErrors: []*pb.SyncUserDataResponse_SyncErrorInfo{
					{
						Text:  "error for card number ",
						Value: testCard.Number,
						Err:   "add card error",
					},
					{
						Text:  "error for pair login/password with prompt ",
						Value: testLoginPwd.Prompt,
						Err:   "add login error",
					},
					{
						Text:  "error for text data with prompt ",
						Value: testTextRecord.Prompt,
						Err:   "add text error",
					},
					{
						Text:  "error for binary data with prompt ",
						Value: testBinaryPb.Prompt,
						Err:   "add bytes error",
					},
				},
				NewLogins:        []*pb.UserLoginPwd{},
				NewCards:         []*pb.UserCard{},
				NewTextRecords:   []*pb.UserTextRecord{},
				NewBinaryRecords: []*pb.UserBinaryRecord{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)

			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.SyncUserData(tt.args.ctx, &pb.SyncUserDataRequest{
				Logins:        tt.inLogins,
				Cards:         tt.inCards,
				TextRecords:   tt.inTexts,
				BinaryRecords: tt.inBinaryes,
				LastSync:      tt.args.lastSync,
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}
		})
	}
}

func TestGetUserCard(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	testTimePrs, err := time.Parse(time.RFC3339, testTime)
	if err != nil {
		fmt.Println("time parse error")
	}

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.Card
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantRes *pb.GetUserCardResponse
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetCard(a.ctx, a.userLogin, a.c.Number).Return(a.c, nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.Card{
					Prompt:    testCard.Prompt,
					Number:    testCard.Number,
					Date:      testCard.Date,
					Code:      testCard.Code,
					Note:      testCard.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserCardResponse{
				Card: &pb.UserCard{
					Prompt:    testCardPb.Prompt,
					Number:    testCardPb.Number,
					Date:      testCardPb.Date,
					Code:      testCardPb.Code,
					Note:      testCardPb.Note,
					TimeStamp: testTime,
				},
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetCard(a.ctx, a.userLogin, a.c.Number).Return(a.c, nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         storage.Card{},
			},
			wantRes: &pb.GetUserCardResponse{
				Card: &pb.UserCard{},
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetCard(a.ctx, a.userLogin, a.c.Number).Return(storage.Card{}, errors.New("error"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.Card{
					Prompt:    testCard.Prompt,
					Number:    testCard.Number,
					Date:      testCard.Date,
					Code:      testCard.Code,
					Note:      testCard.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserCardResponse{
				Card: &pb.UserCard{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.GetUserCard(tt.args.ctx, &pb.GetUserCardRequest{Number: tt.args.c.Number})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}
		})
	}
}

func TestGetUserLogin(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	testTimePrs, err := time.Parse(time.RFC3339, testTime)
	if err != nil {
		fmt.Println("time parse error")
	}

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.LoginPwd
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantRes *pb.GetUserLoginResponse
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login).Return(a.c, nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.LoginPwd{
					Prompt:    testLoginPwd.Prompt,
					Login:     testLoginPwd.Login,
					Pwd:       testLoginPwd.Pwd,
					Note:      testLoginPwd.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserLoginResponse{
				LoginPwd: &pb.UserLoginPwd{
					Prompt:    testLoginPwdPb.Prompt,
					Login:     testLoginPwdPb.Login,
					Pwd:       testLoginPwdPb.Pwd,
					Note:      testLoginPwdPb.Note,
					TimeStamp: testTime,
				},
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login).Return(a.c, nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         storage.LoginPwd{},
			},
			wantRes: &pb.GetUserLoginResponse{
				LoginPwd: &pb.UserLoginPwd{},
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login).Return(storage.LoginPwd{}, errors.New("error"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.LoginPwd{
					Prompt:    testLoginPwd.Prompt,
					Login:     testLoginPwd.Login,
					Pwd:       testLoginPwd.Pwd,
					Note:      testLoginPwd.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserLoginResponse{
				LoginPwd: &pb.UserLoginPwd{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.GetUserLogin(tt.args.ctx, &pb.GetUserLoginRequest{
				Prompt: tt.args.c.Prompt,
				Login:  tt.args.c.Login,
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}
		})
	}
}

func TestGetUserText(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	testTimePrs, err := time.Parse(time.RFC3339, testTime)
	if err != nil {
		fmt.Println("time parse error")
	}

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.TextRecord
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantRes *pb.GetUserTextResponse
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetTextRecord(a.ctx, a.userLogin, a.c.Prompt).Return(a.c, nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.TextRecord{
					Prompt:    testTextRecord.Prompt,
					Data:      testTextRecord.Data,
					Note:      testTextRecord.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserTextResponse{
				TextRecord: &pb.UserTextRecord{
					Prompt:    testTextPb.Prompt,
					Data:      testTextPb.Data,
					Note:      testTextPb.Note,
					TimeStamp: testTime,
				},
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetTextRecord(a.ctx, a.userLogin, a.c.Prompt).Return(a.c, nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         storage.TextRecord{},
			},
			wantRes: &pb.GetUserTextResponse{
				TextRecord: &pb.UserTextRecord{},
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetTextRecord(a.ctx, a.userLogin, a.c.Prompt).Return(storage.TextRecord{}, errors.New("error"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.TextRecord{
					Prompt:    testTextRecord.Prompt,
					Data:      testTextRecord.Data,
					Note:      testTextRecord.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserTextResponse{
				TextRecord: &pb.UserTextRecord{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.GetUserText(tt.args.ctx, &pb.GetUserTextRequest{Prompt: tt.args.c.Prompt})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}
		})
	}
}

func TestGetUserBinary(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	testTimePrs, err := time.Parse(time.RFC3339, testTime)
	if err != nil {
		fmt.Println("time parse error")
	}

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.BinaryRecord
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantRes *pb.GetUserBinaryResponse
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetBinaryRecord(a.ctx, a.userLogin, a.c.Prompt).Return(a.c, nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.BinaryRecord{
					Prompt:    testBinaryRecord.Prompt,
					Data:      testBinaryRecord.Data,
					Note:      testBinaryRecord.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserBinaryResponse{
				BinaryRecord: &pb.UserBinaryRecord{
					Prompt:    testBinaryPb.Prompt,
					Data:      testBinaryPb.Data,
					Note:      testBinaryPb.Note,
					TimeStamp: testTime,
				},
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetBinaryRecord(a.ctx, a.userLogin, a.c.Prompt).Return(a.c, nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         storage.BinaryRecord{},
			},
			wantRes: &pb.GetUserBinaryResponse{
				BinaryRecord: &pb.UserBinaryRecord{},
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetBinaryRecord(a.ctx, a.userLogin, a.c.Prompt).Return(storage.BinaryRecord{}, errors.New("error"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.BinaryRecord{
					Prompt:    testBinaryRecord.Prompt,
					Data:      testBinaryRecord.Data,
					Note:      testBinaryRecord.Note,
					TimeStamp: testTimePrs,
				},
			},
			wantRes: &pb.GetUserBinaryResponse{
				BinaryRecord: &pb.UserBinaryRecord{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := testGRPC.GetUserBinary(tt.args.ctx, &pb.GetUserBinaryRequest{Prompt: tt.args.c.Prompt})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}
		})
	}
}

func TestForceUpdateCard(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.Card
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.Card{
					Prompt:    nil,
					Number:    nil,
					Date:      nil,
					Code:      nil,
					Note:      nil,
					TimeStamp: time.Time{},
				},
				timeSt: testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time stamp test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateCard(a.ctx, a.userLogin, a.c.Prompt, a.c.Number,
					a.c.Date, a.c.Code, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testCard,
				timeSt:    "2006-01T15:04:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.ForceUpdateCard(tt.args.ctx, &pb.ForceUpdateCardRequest{
				Card: &pb.UserCard{
					Prompt:    tt.args.c.Prompt,
					Number:    tt.args.c.Number,
					Date:      tt.args.c.Date,
					Code:      tt.args.c.Code,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestForceUpdateLogin(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.LoginPwd
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.LoginPwd{
					Prompt:    nil,
					Login:     nil,
					Pwd:       nil,
					Note:      nil,
					TimeStamp: time.Time{},
				},
				timeSt: testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time stamp test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateLoginPwd(a.ctx, a.userLogin, a.c.Prompt, a.c.Login, a.c.Pwd, a.c.Note, tp).
					Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testLoginPwd,
				timeSt:    "2006-01-02T154:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.ForceUpdateLoginPwd(tt.args.ctx, &pb.ForceUpdateLoginPwdRequest{
				LoginPwd: &pb.UserLoginPwd{
					Prompt:    tt.args.c.Prompt,
					Login:     tt.args.c.Login,
					Pwd:       tt.args.c.Pwd,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestForceUpdateTextData(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.TextRecord
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c: storage.TextRecord{
					Prompt:    nil,
					Data:      nil,
					Note:      nil,
					TimeStamp: time.Time{},
				},
				timeSt: testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time parse test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateTextRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testTextRecord,
				timeSt:    "2006-01T15:04:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.ForceUpdateTextRecord(tt.args.ctx, &pb.ForceUpdateTextRecordRequest{
				TextRecord: &pb.UserTextRecord{
					Prompt:    tt.args.c.Prompt,
					Data:      tt.args.c.Data,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestForceUpdateBinaryData(t *testing.T) {
	userToken, err := authorizer.BuildToken(testUserLogin, testUserPwd, testCfg.SecretKey)
	if err != nil {
		fmt.Println("build token error")
		return
	}
	ctxWithValue := context.WithValue(context.Background(), authorizer.UserContextKey, userToken)

	type args struct {
		ctx       context.Context
		userLogin string
		c         storage.BinaryRecord
		timeSt    string
	}

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		args    args
		wantErr bool
	}{
		{
			name: "ok test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil)
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: false,
		},
		{
			name: "missing login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).Return(nil).AnyTimes()
			},
			args: args{
				ctx:       context.Background(),
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "empty values test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.EmptyValues, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "exists newer test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(storage.NewStorError(storage.ExistsDataNewerVersion, errors.New("err")))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, a.timeSt)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(errors.New("err"))
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    testTime,
			},
			wantErr: true,
		},
		{
			name: "error time parse test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				tp, err := time.Parse(time.RFC3339, testTime)
				require.NoError(t, err)
				m.EXPECT().ForceUpdateBinaryRecord(a.ctx, a.userLogin, a.c.Prompt, a.c.Data, a.c.Note, tp).
					Return(nil).AnyTimes()
			},
			args: args{
				ctx:       ctxWithValue,
				userLogin: testUserLogin,
				c:         testBinaryRecord,
				timeSt:    "2006-01T15:04:05Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			testGRPC := NewKeeperServer(m, testCfg)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			_, err := testGRPC.ForceUpdateBinaryRecord(tt.args.ctx, &pb.ForceUpdateBinaryRecordRequest{
				BinaryRecord: &pb.UserBinaryRecord{
					Prompt:    tt.args.c.Prompt,
					Data:      tt.args.c.Data,
					Note:      tt.args.c.Note,
					TimeStamp: tt.args.timeSt,
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
