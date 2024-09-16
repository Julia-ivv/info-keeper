package cmdexecutor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cryptor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

// UserCard хранит данные банковской карты.
type UserCard struct {
	Prompt    string
	Number    string
	Date      string
	Code      string
	Note      string
	TimeStamp string
}

// Cards используется для вывода результата пользователю.
type Cards []UserCard

// PrintData используется для вывода результата пользователю.
func (c Cards) PrintData() {
	fmt.Println("CARD")
	for _, v := range c {
		fmt.Println("Prompt: ", v.Prompt)
		fmt.Println("Number: ", v.Number)
		fmt.Println("Date: ", v.Date)
		fmt.Println("Code: ", v.Code)
		fmt.Println("Note: ", v.Note)
		fmt.Println("Time Stamp: ", v.TimeStamp)
	}
}

var addCardExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	err = repo.AddCard(context.Background(), UserLogin, enA.Prompt, enA.CardNumber, enA.CardDate, enA.CardCode,
		enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var updCardExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	err = repo.UpdateCard(context.Background(), UserLogin, enA.Prompt, enA.CardNumber, enA.CardDate, enA.CardCode,
		enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getCardExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enN, err := cryptor.EncryptsString(args.CardNumber)
	if err != nil {
		return nil, err
	}

	c, err := repo.GetCard(context.Background(), UserLogin, enN)
	if err != nil {
		return nil, err
	}

	deC, err := decryptCard(c)
	if err != nil {
		return nil, err
	}

	res := make(Cards, 0, 1)
	res = append(res, deC)

	return res, nil
}

var getCardsExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	cs, err := repo.GetUserCardsAfterTime(context.Background(), UserLogin, time.Now().AddDate(-100, 0, 0).Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	res := make(Cards, 0, len(cs))
	for _, v := range cs {
		c, err := decryptCard(v)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}

	return res, nil
}

var forceAddCardServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enN, err := cryptor.EncryptsString(args.CardNumber)
	if err != nil {
		return nil, err
	}

	c, err := repo.GetCard(context.Background(), UserLogin, enN)
	if err != nil {
		return nil, err
	}

	cPb := cardToPb(c)
	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	_, err = cl.ForceUpdateCard(ctxMd, &pb.ForceUpdateCardRequest{Card: cPb})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getCardServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enN, err := cryptor.EncryptsString(args.CardNumber)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	r, err := cl.GetUserCard(ctxMd, &pb.GetUserCardRequest{Number: enN})
	if err != nil {
		return nil, err
	}
	c := pbToCard(r.Card)
	deC, err := decryptCard(c)
	if err != nil {
		return nil, err
	}

	res := make(Cards, 0, 1)
	res = append(res, UserCard{
		Prompt:    deC.Prompt,
		Number:    deC.Number,
		Date:      deC.Date,
		Code:      deC.Code,
		Note:      deC.Note,
		TimeStamp: deC.TimeStamp,
	})

	return res, nil
}
