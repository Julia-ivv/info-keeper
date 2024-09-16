package cmdexecutor

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/mocks"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
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
		Pwd:       ttArgs.Pwd,
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
		Pwd:        "pwd",
		Text:       "text",
		Binary:     nameTestFile,
	}
	testSyncTime = "2023-01-02T15:04:05Z"
)

func TestExecuteCmd(t *testing.T) {
	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient)
		userCmd string
		args    cmdparser.UserArgs
		wantRes bool
		wantErr bool
		res     DataPrinter
	}{
		{
			name: "ok add card test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().AddCard(context.Background(), "", testCard.Prompt, testCard.Number,
					testCard.Date, testCard.Code, testCard.Note, gomock.Any()).Return(nil)
			},
			userCmd: cmdparser.CmdAddCard,
			args:    ttArgs,
			wantErr: false,
		},
		{
			name: "ok upd card test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().UpdateCard(context.Background(), "", testCard.Prompt, testCard.Number,
					testCard.Date, testCard.Code, testCard.Note, gomock.Any()).Return(nil)
			},
			userCmd: cmdparser.CmdUpdCard,
			args:    ttArgs,
			wantErr: false,
		},
		{
			name: "ok get card test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetCard(context.Background(), "", testCard.Number).Return(testCard, nil)
			},
			userCmd: cmdparser.CmdGetCard,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserCards,
		},
		{
			name: "ok get cards test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetUserCardsAfterTime(context.Background(), "", time.Now().AddDate(-100, 0, 0).Format(time.RFC3339)).
					Return([]storage.Card{testCard}, nil)
			},
			userCmd: cmdparser.CmdGetCards,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserCards,
		},
		{
			name: "ok force add card test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetCard(context.Background(), "", testCard.Number).
					Return(testCard, nil)
				mcli.EXPECT().ForceUpdateCard(ctxMd, &pb.ForceUpdateCardRequest{Card: cardToPb(testCard)}).
					Return(nil, nil)
			},
			userCmd: cmdparser.CmdForceAddCardServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: false,
		},
		{
			name: "ok get server card test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				mcli.EXPECT().GetUserCard(ctxMd, &pb.GetUserCardRequest{Number: testCard.Number}).
					Return(&pb.GetUserCardResponse{Card: cardToPb(testCard)}, nil)
			},
			userCmd: cmdparser.CmdGetCardServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserCards,
		},

		{
			name: "ok get login test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetLoginPwd(context.Background(), "", testLoginPwd.Prompt, testLoginPwd.Login).
					Return(testLoginPwd, nil)
			},
			userCmd: cmdparser.CmdGetLogin,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserLogins,
		},
		{
			name: "ok get logins test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetUserLoginsPwdsAfterTime(context.Background(), "", time.Now().AddDate(-100, 0, 0).Format(time.RFC3339)).
					Return([]storage.LoginPwd{testLoginPwd}, nil)
			},
			userCmd: cmdparser.CmdGetLogins,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserLogins,
		},
		{
			name: "ok force add login test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetLoginPwd(context.Background(), "", testLoginPwd.Prompt, testLoginPwd.Login).
					Return(testLoginPwd, nil)
				mcli.EXPECT().ForceUpdateLoginPwd(ctxMd, &pb.ForceUpdateLoginPwdRequest{LoginPwd: loginToPb(testLoginPwd)}).
					Return(nil, nil)
			},
			userCmd: cmdparser.CmdForceAddLoginServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: false,
		},
		{
			name: "ok get server login test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				mcli.EXPECT().GetUserLogin(ctxMd, &pb.GetUserLoginRequest{Prompt: testLoginPwd.Prompt, Login: testLoginPwd.Login}).
					Return(&pb.GetUserLoginResponse{LoginPwd: loginToPb(testLoginPwd)}, nil)
			},
			userCmd: cmdparser.CmdGetLoginServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserLogins,
		},

		{
			name: "ok add text test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().AddTextRecord(context.Background(), "", testTextRecord.Prompt, testTextRecord.Data,
					testTextRecord.Note, gomock.Any()).Return(nil)
			},
			userCmd: cmdparser.CmdAddText,
			args:    ttArgs,
			wantErr: false,
		},
		{
			name: "ok upd text test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().UpdateTextRecord(context.Background(), "", testTextRecord.Prompt, testTextRecord.Data,
					testTextRecord.Note, gomock.Any()).Return(nil)
			},
			userCmd: cmdparser.CmdUpdText,
			args:    ttArgs,
			wantErr: false,
		},
		{
			name: "ok get text test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetTextRecord(context.Background(), "", testTextRecord.Prompt).Return(testTextRecord, nil)
			},
			userCmd: cmdparser.CmdGetText,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserTexts,
		},
		{
			name: "ok get texts test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetUserTextRecordsAfterTime(context.Background(), "", time.Now().AddDate(-100, 0, 0).Format(time.RFC3339)).
					Return([]storage.TextRecord{testTextRecord}, nil)
			},
			userCmd: cmdparser.CmdGetTexts,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserTexts,
		},
		{
			name: "ok force add text test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetTextRecord(context.Background(), "", testTextRecord.Prompt).
					Return(testTextRecord, nil)
				mcli.EXPECT().ForceUpdateTextRecord(ctxMd, &pb.ForceUpdateTextRecordRequest{TextRecord: textToPb(testTextRecord)}).
					Return(nil, nil)
			},
			userCmd: cmdparser.CmdForceAddTextServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: false,
		},
		{
			name: "ok get server text test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				mcli.EXPECT().GetUserText(ctxMd, &pb.GetUserTextRequest{Prompt: testTextRecord.Prompt}).
					Return(&pb.GetUserTextResponse{TextRecord: textToPb(testTextRecord)}, nil)
			},
			userCmd: cmdparser.CmdGetTextServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserTexts,
		},

		{
			name: "ok add bytes test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().AddBinaryRecord(context.Background(), "", testBinaryRecord.Prompt,
					[]byte{2, 68, 255, 117, 104, 167, 77, 151, 89, 98, 94, 149, 234, 119, 65, 219, 240, 237, 251, 88, 23, 159, 46, 250, 216},
					testBinaryRecord.Note, gomock.Any()).Return(nil)
			},
			userCmd: cmdparser.CmdAddBinary,
			args:    ttArgs,
			wantErr: false,
		},
		{
			name: "ok upd bytes test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().UpdateBinaryRecord(context.Background(), "", testBinaryRecord.Prompt,
					[]byte{2, 68, 255, 117, 104, 167, 77, 151, 89, 98, 94, 149, 234, 119, 65, 219, 240, 237, 251, 88, 23, 159, 46, 250, 216},
					testBinaryRecord.Note, gomock.Any()).Return(nil)
			},
			userCmd: cmdparser.CmdUpdBinary,
			args:    ttArgs,
			wantErr: false,
		},
		{
			name: "ok get bytes test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetBinaryRecord(context.Background(), "", testBinaryRecord.Prompt).Return(testBinaryRecord, nil)
			},
			userCmd: cmdparser.CmdGetBinary,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserBinarys,
		},
		{
			name: "ok get all bytes test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetUserBinaryRecordsAfterTime(context.Background(), "", time.Now().AddDate(-100, 0, 0).Format(time.RFC3339)).
					Return([]storage.BinaryRecord{testBinaryRecord}, nil)
			},
			userCmd: cmdparser.CmdGetBinarys,
			args:    ttArgs,
			wantErr: false,
			wantRes: true,
			res:     testUserBinarys,
		},
		{
			name: "ok force add bytes test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				m.EXPECT().GetBinaryRecord(context.Background(), "", testBinaryRecord.Prompt).
					Return(testBinaryRecord, nil)
				mcli.EXPECT().ForceUpdateBinaryRecord(ctxMd, &pb.ForceUpdateBinaryRecordRequest{BinaryRecord: binaryToPb(testBinaryRecord)}).
					Return(nil, nil)
			},
			userCmd: cmdparser.CmdForceAddBinaryServer,
			args:    ttArgs,
			wantErr: false,
			wantRes: false,
		},
		{
			name: "ok get server bytes test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				mcli.EXPECT().GetUserBinary(ctxMd, &pb.GetUserBinaryRequest{Prompt: testBinaryRecord.Prompt}).
					Return(&pb.GetUserBinaryResponse{BinaryRecord: binaryToPb(testBinaryRecord)}, nil)
			},
			userCmd: cmdparser.CmdGetBinaryServer,
			args:    ttArgs,
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

			ctrlCli := gomock.NewController(t)
			defer ctrlCli.Finish()
			mCli := mocks.NewMockInfoKeeperClient(ctrlCli)

			if tt.prepare != nil {
				tt.prepare(m, mCli)
			}
			res, err := ExecuteCmd(tt.userCmd, tt.args, mCli, m)
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

func TestSynchronization(t *testing.T) {
	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)

	tests := []struct {
		name    string
		prepare func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient)
		wantRes bool
		wantErr bool
		res     DataPrinter
	}{
		{
			name: "ok sync test",
			prepare: func(m *mocks.MockRepositorier, mcli *mocks.MockInfoKeeperClient) {
				gomock.InOrder(
					m.EXPECT().GetLastSyncTime(context.Background(), "").
						Return(testSyncTime, nil),
					m.EXPECT().GetUserCardsAfterTime(context.Background(), "", testSyncTime).
						Return([]storage.Card{testCard}, nil),
					m.EXPECT().GetUserLoginsPwdsAfterTime(context.Background(), "", testSyncTime).
						Return([]storage.LoginPwd{testLoginPwd}, nil),
					m.EXPECT().GetUserTextRecordsAfterTime(context.Background(), "", testSyncTime).
						Return([]storage.TextRecord{testTextRecord}, nil),
					m.EXPECT().GetUserBinaryRecordsAfterTime(context.Background(), "", testSyncTime).
						Return([]storage.BinaryRecord{testBinaryRecord}, nil),
					mcli.EXPECT().SyncUserData(ctxMd, &pb.SyncUserDataRequest{
						Logins:        []*pb.UserLoginPwd{loginToPb(testLoginPwd)},
						Cards:         []*pb.UserCard{cardToPb(testCard)},
						TextRecords:   []*pb.UserTextRecord{textToPb(testTextRecord)},
						BinaryRecords: []*pb.UserBinaryRecord{binaryToPb(testBinaryRecord)},
						LastSync:      testSyncTime,
					}).Return(&pb.SyncUserDataResponse{
						SyncErrors: []*pb.SyncUserDataResponse_SyncErrorInfo{{
							Text:  "text error",
							Value: []byte{20, 89, 224, 162, 229, 20, 169, 198, 23, 48, 193, 238, 14, 23, 152, 188, 173, 160, 95},
							Err:   "error",
						}},
						NewLogins:        []*pb.UserLoginPwd{loginToPb(testLoginPwd)},
						NewCards:         []*pb.UserCard{cardToPb(testCard)},
						NewTextRecords:   []*pb.UserTextRecord{textToPb(testTextRecord)},
						NewBinaryRecords: []*pb.UserBinaryRecord{binaryToPb(testBinaryRecord)},
					}, nil),
					m.EXPECT().AddSyncData(context.Background(), "",
						[]storage.Card{testCard}, []storage.LoginPwd{testLoginPwd},
						[]storage.TextRecord{testTextRecord}, []storage.BinaryRecord{testBinaryRecord}).
						Return(nil),
					m.EXPECT().UpdateLastSyncTime(context.Background(), "", gomock.Any()).
						Return(nil),
				)
			},
			wantErr: false,
			wantRes: true,
			res: SyncErrs{{
				Text:   "text error",
				Value:  "err",
				ErrMsg: "error",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockRepositorier(ctrl)

			ctrlCli := gomock.NewController(t)
			defer ctrlCli.Finish()
			mCli := mocks.NewMockInfoKeeperClient(ctrlCli)

			if tt.prepare != nil {
				tt.prepare(m, mCli)
			}
			resSync, err := synchronization(mCli, m)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.wantRes {
					assert.NotEmpty(t, resSync)
					assert.Equal(t, tt.res, resSync)
				}
			}
		})
	}
}
