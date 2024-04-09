package cmdexecutor

import (
	"context"
	"fmt"
	"time"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cryptor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
	"google.golang.org/grpc/metadata"
)

type UserTextRecord struct {
	Prompt    string
	Data      string
	Note      string
	TimeStamp string
}
type TextRecords []UserTextRecord

func (t TextRecords) PrintData() {
	fmt.Println("TEXT RECORD")
	for _, v := range t {
		fmt.Println("Prompt: ", v.Prompt)
		fmt.Println("Data: ", v.Data)
		fmt.Println("Note: ", v.Note)
		fmt.Println("Time Stamp: ", v.TimeStamp)
	}
}

var addTextExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	err = repo.AddTextRecord(context.Background(), UserLogin, enA.Prompt, enA.Text, enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var updTextExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	err = repo.UpdateTextRecord(context.Background(), UserLogin, enA.Prompt, enA.Text, enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getTextExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}

	t, err := repo.GetTextRecord(context.Background(), UserLogin, enP)
	if err != nil {
		return nil, err
	}

	deT, err := decryptTextRecord(t)
	if err != nil {
		return nil, err
	}

	res := make(TextRecords, 0, 1)
	res = append(res, deT)

	return res, nil
}

var getTextsExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	ts, err := repo.GetUserTextRecordsAfterTime(context.Background(), UserLogin, time.Now().AddDate(-100, 0, 0).Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	res := make(TextRecords, 0, len(ts))
	for _, v := range ts {
		t, err := decryptTextRecord(v)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

var forceAddTextServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}
	t, err := repo.GetTextRecord(context.Background(), UserLogin, enP)
	if err != nil {
		return nil, err
	}

	tPb := textToPb(t)
	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	_, err = cl.ForceUpdateTextRecord(ctxMd, &pb.ForceUpdateTextRecordRequest{TextRecord: tPb})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getTextServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	r, err := cl.GetUserText(ctxMd, &pb.GetUserTextRequest{Prompt: enP})
	if err != nil {
		return nil, err
	}

	t := pbToText(r.TextRecord)
	deT, err := decryptTextRecord(t)
	if err != nil {
		return nil, err
	}

	res := make(TextRecords, 0, 1)
	res = append(res, UserTextRecord{
		Prompt:    deT.Prompt,
		Data:      deT.Data,
		Note:      deT.Note,
		TimeStamp: deT.TimeStamp,
	})

	return res, nil
}
