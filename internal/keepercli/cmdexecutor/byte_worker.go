package cmdexecutor

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cryptor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

type UserBinaryRecord struct {
	Prompt    string
	Data      []byte
	File      string
	Note      string
	TimeStamp string
}

type BinaryRecords []UserBinaryRecord

func (b BinaryRecords) PrintData() {
	fmt.Println("BINARY RECORD")
	for _, v := range b {
		fmt.Println("Prompt: ", v.Prompt)
		fmt.Println("File: ", v.File)
		fmt.Println("Note: ", v.Note)
		fmt.Println("Time Stamp: ", v.TimeStamp)
	}
}

var addBinaryExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	file, err := os.Open(args.Binary)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	args.Binary = ""

	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}
	enData, err := cryptor.EncryptsByte(data)
	if err != nil {
		return nil, err
	}

	err = repo.AddBinaryRecord(context.Background(), UserLogin, enA.Prompt, enData, enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var updBinaryExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	file, err := os.Open(args.Binary)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	args.Binary = ""

	enA, err := encryptArgs(args)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}
	enData, err := cryptor.EncryptsByte(data)
	if err != nil {
		return nil, err
	}

	err = repo.UpdateBinaryRecord(context.Background(), UserLogin, enA.Prompt, enData, enA.Note, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getBinaryExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}
	b, err := repo.GetBinaryRecord(context.Background(), UserLogin, enP)
	if err != nil {
		return nil, err
	}

	deB, err := decryptBinaryRecord(b)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(args.Prompt, deB.Data, 0666)
	if err != nil {
		return nil, err
	}
	deB.File = args.Prompt
	deB.Data = nil

	res := make(BinaryRecords, 0, 1)
	res = append(res, deB)

	return res, nil
}

var getBinarysExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	bs, err := repo.GetUserBinaryRecordsAfterTime(context.Background(), UserLogin, time.Now().AddDate(-100, 0, 0).Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	res := make(BinaryRecords, 0, len(bs))
	for _, v := range bs {
		b, err := decryptBinaryRecord(v)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(args.Prompt, b.Data, 0666)
		if err != nil {
			return nil, err
		}
		b.Data = nil
		b.File = b.Prompt
		res = append(res, b)
	}

	return res, nil
}

var forceAddBinaryServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}
	b, err := repo.GetBinaryRecord(context.Background(), UserLogin, enP)
	if err != nil {
		return nil, err
	}

	bPb := binaryToPb(b)
	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	_, err = cl.ForceUpdateBinaryRecord(ctxMd, &pb.ForceUpdateBinaryRecordRequest{BinaryRecord: bPb})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getBinaryServerExec = func(args cmdparser.UserArgs, cl pb.InfoKeeperClient, repo storage.Repositorier) (DataPrinter, error) {
	enP, err := cryptor.EncryptsString(args.Prompt)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{authorizer.AccessToken: UserToken})
	ctxMd := metadata.NewOutgoingContext(context.Background(), md)
	r, err := cl.GetUserBinary(ctxMd, &pb.GetUserBinaryRequest{Prompt: enP})
	if err != nil {
		return nil, err
	}

	b := pbToBinary(r.BinaryRecord)
	deB, err := decryptBinaryRecord(b)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(args.Prompt, deB.Data, 0666)
	if err != nil {
		return nil, err
	}
	deB.File = args.Prompt
	deB.Data = nil

	res := make(BinaryRecords, 0, 1)
	res = append(res, UserBinaryRecord{
		Prompt:    deB.Prompt,
		Data:      deB.Data,
		File:      deB.File,
		Note:      deB.Note,
		TimeStamp: deB.TimeStamp,
	})

	return res, nil
}
