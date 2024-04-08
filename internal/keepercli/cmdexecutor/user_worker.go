package cmdexecutor

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.com/david_mbuvi/go_asterisks"
	"google.golang.org/grpc/metadata"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cryptor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

var UserToken string
var UserLogin string
var UserPath string

var regExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	fmt.Print("Enter your password: ")
	password, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter key: ")
	cryptor.UserKey, err = go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	resp, err := cl.AddUser(context.Background(), &pb.AddUserRequest{Login: args.AuthLogin, Pwd: string(password)})
	if err != nil {
		return nil, err
	}
	UserToken = resp.GetToken()

	err = repo.RegUser(context.Background(), args.AuthLogin, string(password))
	if err != nil {
		return nil, err
	}

	err = repo.AuthUser(context.Background(), args.AuthLogin, string(password))
	if err != nil {
		return nil, err
	}
	UserLogin = args.AuthLogin
	UserPath = "./" + UserLogin + "/"

	return nil, nil
}

type SyncErr struct {
	Text   string
	Value  string
	ErrMsg string
}

type SyncErrs []SyncErr

func (s SyncErrs) PrintData() {
	for _, v := range s {
		fmt.Println("SYNC ERRORS")
		fmt.Printf("%s %s: %s\n", v.Text, v.Value, v.ErrMsg)
	}
}

func synchronization(cl pb.InfoKeeperClient, repo storage.Repositorier) (SyncErrs, error) {
	lSync, err := repo.GetLastSyncTime(context.Background(), UserLogin)
	if err != nil {
		return nil, err
	}
	cs, err := repo.GetUserCardsAfterTime(context.Background(), UserLogin, lSync)
	if err != nil {
		return nil, err
	}
	ls, err := repo.GetUserLoginsPwdsAfterTime(context.Background(), UserLogin, lSync)
	if err != nil {
		return nil, err
	}
	ts, err := repo.GetUserTextRecordsAfterTime(context.Background(), UserLogin, lSync)
	if err != nil {
		return nil, err
	}
	bs, err := repo.GetUserBinaryRecordsAfterTime(context.Background(), UserLogin, lSync)
	if err != nil {
		return nil, err
	}

	pbC := cardsToPb(cs)
	pbL := loginsToPb(ls)
	pbT := textsToPb(ts)
	pbB := binarysToPb(bs)

	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	resSync, err := cl.SyncUserData(ctxMd, &pb.SyncUserDataRequest{
		Logins:        pbL,
		Cards:         pbC,
		TextRecords:   pbT,
		BinaryRecords: pbB,
		LastSync:      lSync,
	})
	if err != nil {
		return nil, err
	}

	newCs := pbToCards(resSync.GetNewCards())
	newLs := pbToLogins(resSync.GetNewLogins())
	newTs := pbToTexts(resSync.GetNewTextRecords())
	newBs := pbToBinarys(resSync.GetNewBinaryRecords())

	err = repo.AddSyncData(context.Background(), UserLogin, newCs, newLs, newTs, newBs)
	if err != nil {
		return nil, err
	}

	err = repo.UpdateLastSyncTime(context.Background(), UserLogin, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	r := make(SyncErrs, 0, len(resSync.SyncErrors))
	for _, v := range resSync.SyncErrors {
		val, err := cryptor.Decrypts(v.Value)
		if err != nil {
			val = "decryption error"
		}
		r = append(r, SyncErr{
			Text:   v.Text,
			Value:  val,
			ErrMsg: v.Err,
		})
	}

	return r, nil
}

var authExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	fmt.Print("Enter your password: ")
	password, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter key: ")
	cryptor.UserKey, err = go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	resp, err := cl.AuthUser(context.Background(), &pb.AuthUserRequest{Login: UserLogin, Pwd: string(password)})
	if err != nil {
		return nil, err
	}

	err = repo.AuthUser(context.Background(), args.AuthLogin, string(password))
	if err != nil {
		return nil, err
	}

	UserLogin = args.AuthLogin
	UserToken = resp.GetToken()
	UserPath = "./" + UserLogin + "/"

	return synchronization(cl, repo)
}

var exitExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	r, err := synchronization(cl, repo)
	if err != nil {
		fmt.Println(err)
	}
	r.PrintData()
	os.Exit(0)
	return nil, nil
}
