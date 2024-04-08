package cmdexecutor

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cryptor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
	"gitlab.com/david_mbuvi/go_asterisks"
	"google.golang.org/grpc/metadata"
)

type UserLoginPwd struct {
	Prompt    string
	Login     string
	Pwd       string
	Note      string
	TimeStamp string
}
type LoginPwds []UserLoginPwd

func (l LoginPwds) PrintData() {
	fmt.Println("LOGIN PWD")
	for _, v := range l {
		fmt.Println("Prompt: ", v.Prompt)
		fmt.Println("Login: ", v.Login)
		fmt.Println("Pwd: ", v.Pwd)
		fmt.Println("Note: ", v.Note)
		fmt.Println("Time Stamp: ", v.TimeStamp)
	}
}

var addLoginExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter password: ")
	pwd, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	err = repo.AddLoginPwd(context.Background(), UserLogin, enA.Prompt, enA.Login, pwd, enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var updLoginExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter password: ")
	pwd, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	err = repo.UpdateLoginPwd(context.Background(), UserLogin, enA.Prompt, enA.Login, pwd, enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getLoginExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}
	enL, err := cryptor.EncryptsString(args.Login)
	if err != nil {
		return nil, err
	}

	l, err := repo.GetLoginPwd(context.Background(), UserLogin, enP, enL)
	if err != nil {
		return nil, err
	}

	deLP, err := decryptLoginPwd(l)
	if err != nil {
		return nil, err
	}

	res := make(LoginPwds, 0, 1)
	res = append(res, deLP)

	return res, nil
}

var getLoginsExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	ls, err := repo.GetUserLoginsPwdsAfterTime(context.Background(), UserLogin, time.Now().AddDate(100, 0, 0).Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	res := make(LoginPwds, 0, len(ls))
	for _, v := range ls {
		l, err := decryptLoginPwd(v)
		if err != nil {
			return nil, err
		}
		res = append(res, l)
	}

	return res, nil
}

var forceAddLoginServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}
	enL, err := cryptor.EncryptsString(args.Login)
	if err != nil {
		return nil, err
	}

	l, err := repo.GetLoginPwd(context.Background(), UserLogin, enP, enL)
	if err != nil {
		return nil, err
	}

	lPb := loginToPb(l)
	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	_, err = cl.ForceUpdateLoginPwd(ctxMd, &pb.ForceUpdateLoginPwdRequest{LoginPwd: lPb})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getLoginServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}
	enL, err := cryptor.EncryptsString(args.Login)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	r, err := cl.GetUserLogin(ctxMd, &pb.GetUserLoginRequest{Prompt: enP, Login: enL})
	if err != nil {
		return nil, err
	}

	l := pbToLogin(r.LoginPwd)
	deL, err := decryptLoginPwd(l)
	if err != nil {
		return nil, err
	}

	res := make(LoginPwds, 0, 1)
	res = append(res, UserLoginPwd{
		Prompt:    deL.Prompt,
		Login:     deL.Login,
		Pwd:       deL.Pwd,
		Note:      deL.Note,
		TimeStamp: deL.TimeStamp,
	})

	return res, nil
}
