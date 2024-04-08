package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdexecutor"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/cmdparser"
	cliConfig "github.com/Julia-ivv/info-keeper.git/internal/keepercli/config"
	"github.com/Julia-ivv/info-keeper.git/internal/keepercli/storage"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

func main() {
	cfg := cliConfig.NewConfig()
	logger.ZapSugar = logger.NewLogger()
	logger.ZapSugar.Infow("Starting gRPC client", "port", cfg.GRPC)
	logger.ZapSugar.Infow("Database", "path", cfg.DBURI)

	if err := keyboard.Open(); err != nil {
		logger.ZapSugar.Fatal(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	repo, err := storage.NewStorage(*cfg)
	if err != nil {
		logger.ZapSugar.Fatal(err)
	}
	defer repo.Close()

	conn, err := grpc.NewClient(cfg.GRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cl := pb.NewInfoKeeperClient(conn)

	for {
		fmt.Println("Enter command: ")
		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			logger.ZapSugar.Infoln("can`t read command.", err)
		}
		userCmd, userArgs, err := cmdparser.ParseUserCmd(userInput)
		if err != nil {
			logger.ZapSugar.Infoln("can`t parse command ", userInput, err)
		} else if userCmd != "" {
			res, err := cmdexecutor.ExecuteCmd(userCmd, userArgs, cl, repo)
			if err != nil {
				logger.ZapSugar.Infoln("can`t execute command ", userCmd, err)
			}
			res.PrintData()
		}
	}
}
