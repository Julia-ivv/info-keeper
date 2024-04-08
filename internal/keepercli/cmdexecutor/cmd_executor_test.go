package cmdexecutor

import (
	"context"
	"testing"
	"time"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/mocks"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testUserCards = Cards{{
		Prompt:    ttArgs.Prompt,
		Number:    ttArgs.CardNumber,
		Date:      ttArgs.CardDate,
		Code:      ttArgs.CardCode,
		Note:      ttArgs.Note,
		TimeStamp: testTime,
	}}
	testUserLogins = LoginPwds{{
		Prompt:    ttArgs.Prompt,
		Login:     ttArgs.Login,
		Pwd:       "pwd",
		Note:      ttArgs.Note,
		TimeStamp: testTime,
	}}
	testUserTexts = TextRecords{{
		Prompt:    ttArgs.Prompt,
		Data:      ttArgs.Text,
		Note:      ttArgs.Note,
		TimeStamp: testTime,
	}}
	testUserBinarys = BinaryRecords{{
		Prompt:    ttArgs.Prompt,
		Data:      []byte{45, 46},
		Note:      ttArgs.Note,
		TimeStamp: testTime,
	}}
	nameTestFile = "test_file"
	ttArgs       = cmdparser.UserArgs{
		Prompt:     "prompt",
		Note:       "note",
		CardNumber: "123",
		CardDate:   "12/24",
		CardCode:   "555",
		Login:      "login",
		Text:       "text",
		Binary:     nameTestFile,
	}
)

func TestExecuteCmd(t *testing.T) {

	type args struct {
		u cmdparser.UserArgs
	}
	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, a args)
		userCmd string
		args    args
		wantRes bool
		wantErr bool
		res     DataPrinter
	}{
		{
			name: "ok add card test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().AddCard(context.Background(), "", testCard.Prompt, testCard.Number,
					testCard.Date, testCard.Code, testCard.Note, time.Now().Format(time.RFC3339)).Return(nil)
			},
			userCmd: cmdparser.CmdAddCard,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
		},
		{
			name: "ok upd card test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().UpdateCard(context.Background(), "", testCard.Prompt, testCard.Number,
					testCard.Date, testCard.Code, testCard.Note, time.Now().Format(time.RFC3339)).Return(nil)
			},
			userCmd: cmdparser.CmdUpdCard,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
		},
		{
			name: "ok get card test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetCard(context.Background(), "", testCard.Number).Return(testCard, nil)
			},
			userCmd: cmdparser.CmdGetCard,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserCards,
		},
		{
			name: "ok get cards test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetUserCardsAfterTime(context.Background(), "", time.Now().AddDate(100, 0, 0).Format(time.RFC3339)).
					Return([]storage.Card{testCard}, nil)
			},
			userCmd: cmdparser.CmdGetCards,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserCards,
		},

		{
			name: "ok get login test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetLoginPwd(context.Background(), "", testLoginPwd.Prompt, testLoginPwd.Login).
					Return(testLoginPwd, nil)
			},
			userCmd: cmdparser.CmdGetLogin,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserLogins,
		},
		{
			name: "ok get logins test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetUserLoginsPwdsAfterTime(context.Background(), "", time.Now().AddDate(100, 0, 0).Format(time.RFC3339)).
					Return([]storage.LoginPwd{testLoginPwd}, nil)
			},
			userCmd: cmdparser.CmdGetLogins,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserLogins,
		},

		{
			name: "ok add text test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().AddTextRecord(context.Background(), "", testTextRecord.Prompt, testTextRecord.Data,
					testTextRecord.Note, time.Now().Format(time.RFC3339)).Return(nil)
			},
			userCmd: cmdparser.CmdAddText,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
		},
		{
			name: "ok upd text test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().UpdateTextRecord(context.Background(), "", testTextRecord.Prompt, testTextRecord.Data,
					testTextRecord.Note, time.Now().Format(time.RFC3339)).Return(nil)
			},
			userCmd: cmdparser.CmdUpdText,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
		},
		{
			name: "ok get text test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetTextRecord(context.Background(), "", testTextRecord.Prompt).Return(testTextRecord, nil)
			},
			userCmd: cmdparser.CmdGetText,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserTexts,
		},
		{
			name: "ok get texts test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetUserTextRecordsAfterTime(context.Background(), "", time.Now().AddDate(100, 0, 0).Format(time.RFC3339)).
					Return([]storage.TextRecord{testTextRecord}, nil)
			},
			userCmd: cmdparser.CmdGetTexts,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserTexts,
		},

		{
			name: "ok add bytes test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().AddBinaryRecord(context.Background(), "", testBinaryRecord.Prompt,
					[]byte{2, 68, 255, 117, 104, 167, 77, 151, 89, 98, 94, 149, 234, 119, 65, 219, 240, 237, 251, 88, 23, 159, 46, 250, 216},
					testBinaryRecord.Note, time.Now().Format(time.RFC3339)).Return(nil)
			},
			userCmd: cmdparser.CmdAddBinary,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
		},
		{
			name: "ok upd bytes test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().UpdateBinaryRecord(context.Background(), "", testBinaryRecord.Prompt,
					[]byte{2, 68, 255, 117, 104, 167, 77, 151, 89, 98, 94, 149, 234, 119, 65, 219, 240, 237, 251, 88, 23, 159, 46, 250, 216},
					testBinaryRecord.Note, time.Now().Format(time.RFC3339)).Return(nil)
			},
			userCmd: cmdparser.CmdUpdBinary,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
		},
		{
			name: "ok get bytes test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetBinaryRecord(context.Background(), "", testBinaryRecord.Prompt).Return(testBinaryRecord, nil)
			},
			userCmd: cmdparser.CmdGetBinary,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserBinarys,
		},
		{
			name: "ok get all bytes test",
			prepare: func(m *mocks.MockRepositorier, a args) {
				m.EXPECT().GetUserBinaryRecordsAfterTime(context.Background(), "", time.Now().AddDate(100, 0, 0).Format(time.RFC3339)).
					Return([]storage.BinaryRecord{testBinaryRecord}, nil)
			},
			userCmd: cmdparser.CmdGetBinarys,
			args: args{
				u: ttArgs,
			},
			wantErr: false,
			wantRes: true,
			res:     testUserBinarys,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockRepositorier(ctrl)
			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}
			res, err := ExecuteCmd(tt.userCmd, tt.args.u, nil, m)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.wantRes {
					assert.NotEmpty(t, res)
				}
			}
		})
	}
}
