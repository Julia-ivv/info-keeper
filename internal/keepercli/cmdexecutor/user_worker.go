package cmdexecutor

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/david_mbuvi/go_asterisks"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

var UserToken string
var UserKey []byte
var UserLogin string

var regExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (res []CmdResult, err error) {
	fmt.Print("Enter your password: ")
	password, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter key: ")
	UserKey, err = go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	resp, err := cl.AddUser(context.Background(), &pb.AddUserRequest{Login: args.UserLogin, Pwd: string(password)})
	if err != nil {
		return nil, err
	}
	UserToken = resp.GetToken()

	err = repo.RegUser(context.Background(), args.UserLogin, string(password))
	if err != nil {
		return nil, err
	}

	err = repo.AuthUser(context.Background(), args.UserLogin, string(password))
	if err != nil {
		return nil, err
	}
	UserLogin = args.UserLogin

	return nil, nil
}
var authExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (res []CmdResult, err error) {
	// here must be sync
	fmt.Print("Enter your password: ")
	password, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter key: ")
	UserKey, err = go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	resp, err := cl.AuthUser(context.Background(), &pb.AuthUserRequest{Login: UserLogin, Pwd: string(password)})
	if err != nil {
		return nil, err
	}

	err = repo.AuthUser(context.Background(), args.UserLogin, string(password))
	if err != nil {
		return nil, err
	}

	UserLogin = args.UserLogin
	UserToken = resp.GetToken()

	// resSync, err := cl.SyncUserData(context.Background(), &pb.SyncUserDataRequest{
	// 	Logins:        []*pb.UserLoginPwd{},
	// 	Cards:         []*pb.UserCard{},
	// 	TextRecords:   []*pb.UserTextRecord{},
	// 	BinaryRecords: []*pb.UserBinaryRecord{},
	// 	LastSync:      "",
	// })
	return nil, nil
}
var exitExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (res []CmdResult, err error) {
	// add sync
	os.Exit(0)
	return nil, nil
}
