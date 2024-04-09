// Пакет grpcserver реализует gRPC-методы.
package grpcserver

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keeper/config"
	"github.com/Julia-ivv/info-keeper.git/internal/keeper/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

// ShortenerServer хранит репозиторий и настройки приложения.
type KeeperGRPCServer struct {
	pb.UnimplementedInfoKeeperServer
	stor storage.Repositorier
	cfg  config.Flags
}

// NewShortenerServer создает объект с репозиторием и настройками для gRPC-методов.
func NewKeeperServer(stor storage.Repositorier, cfg config.Flags) *KeeperGRPCServer {
	k := &KeeperGRPCServer{}
	k.stor = stor
	k.cfg = cfg
	return k
}

// AddUser реализует добавление и аутентификацию пользователя в БД.
func (ks *KeeperGRPCServer) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	if in.GetLogin() == "" || in.GetPwd() == "" {
		return nil, status.Error(codes.DataLoss, "empty login or password")
	}

	err := ks.stor.RegUser(ctx, in.GetLogin(), in.GetPwd())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, status.Error(codes.AlreadyExists, "user with login "+in.Login+" already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = ks.stor.AuthUser(ctx, in.GetLogin(), in.GetPwd())
	if err != nil {
		var authErr *authorizer.AuthErr
		if (errors.As(err, &authErr)) && (authErr.ErrType == authorizer.InvalidHash) {
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenString, err := authorizer.BuildToken(in.Login, in.Pwd, ks.cfg.SecretKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddUserResponse{Token: tokenString}, nil
}

// AuthUser - аутентификация существующего пользователя.
func (ks *KeeperGRPCServer) AuthUser(ctx context.Context, in *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
	if in.GetLogin() == "" || in.GetPwd() == "" {
		return nil, status.Error(codes.DataLoss, "empty login or password")
	}

	err := ks.stor.AuthUser(ctx, in.GetLogin(), in.GetPwd())
	if err != nil {
		var authErr *authorizer.AuthErr
		if (errors.As(err, &authErr)) && (authErr.ErrType == authorizer.InvalidHash) {
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenString, err := authorizer.BuildToken(in.GetLogin(), in.GetPwd(), ks.cfg.SecretKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AuthUserResponse{Token: tokenString}, nil
}

// AddCard реализует добавление информации о банковской карте.
func (ks *KeeperGRPCServer) AddCard(ctx context.Context, in *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetCard() == nil {
		return nil, status.Error(codes.DataLoss, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.Card.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = ks.stor.AddCard(ctx, userLogin, in.Card.GetPrompt(), in.Card.GetNumber(), in.Card.GetDate(),
		in.Card.GetCode(), in.Card.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.As(err, &addErr) && addErr.ErrType == storage.ExistsDataNewerVersion {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// AddLogin реализует добавление пары логин - пароль в БД.
func (ks *KeeperGRPCServer) AddLogin(ctx context.Context, in *pb.AddLoginRequest) (*pb.AddLoginResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetLoginPwd() == nil {
		return nil, status.Error(codes.DataLoss, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.LoginPwd.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = ks.stor.AddLoginPwd(ctx, userLogin, in.LoginPwd.GetPrompt(), in.LoginPwd.GetLogin(), in.LoginPwd.GetPwd(), in.LoginPwd.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.As(err, &addErr) && addErr.ErrType == storage.ExistsDataNewerVersion {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// AddTextData - раелизует добавление текстовой информации.
func (ks *KeeperGRPCServer) AddTextData(ctx context.Context, in *pb.AddTextDataRequest) (*pb.AddTextDataResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetTextRecord() == nil {
		return nil, status.Error(codes.DataLoss, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.TextRecord.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = ks.stor.AddTextRecord(ctx, userLogin, in.TextRecord.GetPrompt(), in.TextRecord.GetData(), in.TextRecord.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.As(err, &addErr) && addErr.ErrType == storage.ExistsDataNewerVersion {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// AddBinaryData реализует добавление бинарной информации в БД.
func (ks *KeeperGRPCServer) AddBinaryData(ctx context.Context, in *pb.AddBinaryDataRequest) (*pb.AddBinaryDataResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetBinaryRecord() == nil {
		return nil, status.Error(codes.DataLoss, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.BinaryRecord.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = ks.stor.AddBinaryRecord(ctx, userLogin, in.BinaryRecord.GetPrompt(), in.BinaryRecord.GetData(), in.BinaryRecord.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.As(err, &addErr) && addErr.ErrType == storage.ExistsDataNewerVersion {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// SyncErrInfo - содержит информацию об ошибках при синхронизации.
type SyncErrInfo struct {
	Text  string
	Value []byte
	Err   string
}

// SyncUserData выполняет синхронизацию данных между сервером и клиентом.
func (ks *KeeperGRPCServer) SyncUserData(ctx context.Context, in *pb.SyncUserDataRequest) (*pb.SyncUserDataResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	lastSync, err := time.Parse(time.RFC3339, string(in.GetLastSync()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	newCards, err := ks.stor.GetUserCardsAfterTime(ctx, userLogin, lastSync)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newLogins, err := ks.stor.GetUserLoginsPwdsAfterTime(ctx, userLogin, lastSync)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newTextRecords, err := ks.stor.GetUserTextRecordsAfterTime(ctx, userLogin, lastSync)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newBinaryRecords, err := ks.stor.GetUserBinaryRecordsAfterTime(ctx, userLogin, lastSync)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var respErrors []SyncErrInfo

	if in.GetCards() != nil {
		for _, v := range in.GetCards() {
			timeStamp, err := time.Parse(time.RFC3339, v.GetTimeStamp())
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for card number ",
					Value: v.GetNumber(),
					Err:   err.Error(),
				})
				continue
			}
			err = ks.stor.AddCard(ctx, userLogin, v.GetPrompt(), v.GetNumber(), v.GetDate(), v.GetCode(), v.GetNote(), timeStamp)
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for card number ",
					Value: v.GetNumber(),
					Err:   err.Error(),
				})
			}
			newCards = slices.DeleteFunc(newCards, func(c storage.Card) bool {
				return slices.Compare(c.Number, v.GetNumber()) == 0
			})
		}
	}

	if in.GetLogins() != nil {
		for _, v := range in.GetLogins() {
			timeStamp, err := time.Parse(time.RFC3339, v.GetTimeStamp())
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for pair login/password with prompt ",
					Value: v.GetPrompt(),
					Err:   err.Error(),
				})
				continue
			}
			err = ks.stor.AddLoginPwd(ctx, userLogin, v.GetPrompt(), v.GetLogin(), v.GetPwd(), v.GetNote(), timeStamp)
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for pair login/password with prompt ",
					Value: v.GetPrompt(),
					Err:   err.Error(),
				})
			}
			newLogins = slices.DeleteFunc(newLogins, func(l storage.LoginPwd) bool {
				return slices.Compare(l.Prompt, v.GetPrompt()) == 0 && slices.Compare(l.Login, v.GetLogin()) == 0
			})
		}
	}

	if in.GetTextRecords() != nil {
		for _, v := range in.GetTextRecords() {
			timeStamp, err := time.Parse(time.RFC3339, v.GetTimeStamp())
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for text data with prompt ",
					Value: v.GetPrompt(),
					Err:   err.Error(),
				})
				continue
			}
			err = ks.stor.AddTextRecord(ctx, userLogin, v.GetPrompt(), v.GetData(), v.GetNote(), timeStamp)
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for text data with prompt ",
					Value: v.GetPrompt(),
					Err:   err.Error(),
				})
			}
			newTextRecords = slices.DeleteFunc(newTextRecords, func(t storage.TextRecord) bool {
				return slices.Compare(t.Prompt, v.GetPrompt()) == 0
			})
		}
	}

	if in.GetBinaryRecords() != nil {
		for _, v := range in.GetBinaryRecords() {
			timeStamp, err := time.Parse(time.RFC3339, v.GetTimeStamp())
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for binary data with prompt ",
					Value: v.GetPrompt(),
					Err:   err.Error(),
				})
				continue
			}
			err = ks.stor.AddBinaryRecord(ctx, userLogin, v.GetPrompt(), v.GetData(), v.GetNote(), timeStamp)
			if err != nil {
				respErrors = append(respErrors, SyncErrInfo{
					Text:  "error for binary data with prompt ",
					Value: v.GetPrompt(),
					Err:   err.Error(),
				})
			}
			newBinaryRecords = slices.DeleteFunc(newBinaryRecords, func(b storage.BinaryRecord) bool {
				return slices.Compare(b.Prompt, v.GetPrompt()) == 0
			})
		}
	}

	respCards := make([]*pb.UserCard, 0, len(newCards))
	for _, v := range newCards {
		respCards = append(respCards, &pb.UserCard{
			Prompt:    v.Prompt,
			Number:    v.Number,
			Date:      v.Date,
			Code:      v.Code,
			Note:      v.Note,
			TimeStamp: v.TimeStamp.Format(time.RFC3339),
		})
	}

	respLogins := make([]*pb.UserLoginPwd, 0, len(newLogins))
	for _, v := range newLogins {
		respLogins = append(respLogins, &pb.UserLoginPwd{
			Prompt:    v.Prompt,
			Login:     v.Login,
			Pwd:       v.Pwd,
			Note:      v.Note,
			TimeStamp: v.TimeStamp.Format(time.RFC3339),
		})
	}

	respText := make([]*pb.UserTextRecord, 0, len(newTextRecords))
	for _, v := range newTextRecords {
		respText = append(respText, &pb.UserTextRecord{
			Prompt:    v.Prompt,
			Data:      v.Data,
			Note:      v.Note,
			TimeStamp: v.TimeStamp.Format(time.RFC3339),
		})
	}

	respBinary := make([]*pb.UserBinaryRecord, 0, len(newBinaryRecords))
	for _, v := range newBinaryRecords {
		respBinary = append(respBinary, &pb.UserBinaryRecord{
			Prompt:    v.Prompt,
			Data:      v.Data,
			Note:      v.Note,
			TimeStamp: v.TimeStamp.Format(time.RFC3339),
		})
	}

	errInfo := make([]*pb.SyncUserDataResponse_SyncErrorInfo, 0, len(respErrors))
	for _, v := range respErrors {
		errInfo = append(errInfo, &pb.SyncUserDataResponse_SyncErrorInfo{
			Text:  v.Text,
			Value: v.Value,
			Err:   v.Err,
		})
	}

	return &pb.SyncUserDataResponse{
		SyncErrors:       errInfo,
		NewLogins:        respLogins,
		NewCards:         respCards,
		NewTextRecords:   respText,
		NewBinaryRecords: respBinary,
	}, nil
}

// GetUserCard реализует получение информации о банковской карте.
func (ks *KeeperGRPCServer) GetUserCard(ctx context.Context, in *pb.GetUserCardRequest) (*pb.GetUserCardResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	card, err := ks.stor.GetCard(ctx, userLogin, in.GetNumber())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserCardResponse{
		Card: &pb.UserCard{
			Prompt:    card.Prompt,
			Number:    card.Number,
			Date:      card.Date,
			Code:      card.Code,
			Note:      card.Note,
			TimeStamp: card.TimeStamp.Format(time.RFC3339),
		},
	}, nil
}

// GetUserLogin - получает информацию о паре логин-пароль из БД.
func (ks *KeeperGRPCServer) GetUserLogin(ctx context.Context, in *pb.GetUserLoginRequest) (*pb.GetUserLoginResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	lg, err := ks.stor.GetLoginPwd(ctx, userLogin, in.GetPrompt(), in.GetLogin())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserLoginResponse{
		LoginPwd: &pb.UserLoginPwd{
			Prompt:    lg.Prompt,
			Login:     lg.Login,
			Pwd:       lg.Pwd,
			Note:      lg.Note,
			TimeStamp: lg.TimeStamp.Format(time.RFC3339),
		},
	}, nil
}

// GetUserText - получает текстовую информацию из БД.
func (ks *KeeperGRPCServer) GetUserText(ctx context.Context, in *pb.GetUserTextRequest) (*pb.GetUserTextResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	tr, err := ks.stor.GetTextRecord(ctx, userLogin, in.GetPrompt())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserTextResponse{
		TextRecord: &pb.UserTextRecord{
			Prompt:    tr.Prompt,
			Data:      tr.Data,
			Note:      tr.Note,
			TimeStamp: tr.TimeStamp.Format(time.RFC3339),
		},
	}, nil
}

// GetUserBinary - получает бинарную информацию из БД.
func (ks *KeeperGRPCServer) GetUserBinary(ctx context.Context, in *pb.GetUserBinaryRequest) (*pb.GetUserBinaryResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	br, err := ks.stor.GetBinaryRecord(ctx, userLogin, in.GetPrompt())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserBinaryResponse{
		BinaryRecord: &pb.UserBinaryRecord{
			Prompt:    br.Prompt,
			Data:      br.Data,
			Note:      br.Note,
			TimeStamp: br.TimeStamp.Format(time.RFC3339),
		},
	}, nil
}

// ForceUpdateCard - обновляет информацию о банковской карте.
func (ks *KeeperGRPCServer) ForceUpdateCard(ctx context.Context, in *pb.ForceUpdateCardRequest) (*pb.ForceUpdateCardResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetCard() == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.Card.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = ks.stor.ForceUpdateCard(ctx, userLogin, in.Card.GetPrompt(), in.Card.GetNumber(), in.Card.GetDate(),
		in.Card.GetCode(), in.Card.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// ForceUpdateLoginPwd - обновляет информацию о паре логин-пароль.
func (ks *KeeperGRPCServer) ForceUpdateLoginPwd(ctx context.Context, in *pb.ForceUpdateLoginPwdRequest) (*pb.ForceUpdateLoginPwdResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetLoginPwd() == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.LoginPwd.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = ks.stor.ForceUpdateLoginPwd(ctx, userLogin, in.LoginPwd.GetPrompt(), in.LoginPwd.GetLogin(), in.LoginPwd.GetPwd(),
		in.LoginPwd.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// ForceUpdateTextRecord - обновляет текстовую информацию.
func (ks *KeeperGRPCServer) ForceUpdateTextRecord(ctx context.Context, in *pb.ForceUpdateTextRecordRequest) (*pb.ForceUpdateTextRecordResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetTextRecord() == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.TextRecord.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = ks.stor.ForceUpdateTextRecord(ctx, userLogin, in.TextRecord.GetPrompt(), in.TextRecord.GetData(), in.TextRecord.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// ForceUpdateBinaryRecord - обновляет бинарную информацию в БД.
func (ks *KeeperGRPCServer) ForceUpdateBinaryRecord(ctx context.Context, in *pb.ForceUpdateBinaryRecordRequest) (*pb.ForceUpdateBinaryRecordResponse, error) {
	v := ctx.Value(authorizer.UserContextKey)
	if v == nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userToken := v.(string)

	userLogin, err := authorizer.GetUserDataFromToken(userToken, ks.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if in.GetBinaryRecord() == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	timeStamp, err := time.Parse(time.RFC3339, in.BinaryRecord.GetTimeStamp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = ks.stor.ForceUpdateBinaryRecord(ctx, userLogin, in.BinaryRecord.GetPrompt(), in.BinaryRecord.GetData(), in.BinaryRecord.GetNote(), timeStamp)
	if err != nil {
		var addErr *storage.StorErr
		if errors.As(err, &addErr) && addErr.ErrType == storage.EmptyValues {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
