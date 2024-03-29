package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/Julia-ivv/info-keeper.git/internal/config"
	"github.com/Julia-ivv/info-keeper.git/internal/grpcserver"
	"github.com/Julia-ivv/info-keeper.git/internal/interceptors"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
	"github.com/Julia-ivv/info-keeper.git/internal/storage"
	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

func main() {
	cfg := config.NewConfig()

	logger.ZapSugar = logger.NewLogger()
	logger.ZapSugar.Infow("Starting gRPC server", "port", cfg.GRPC)
	logger.ZapSugar.Infow("flags", "db dsn", cfg.DBDSN)

	repo, err := storage.NewStorage(*cfg)
	if err != nil {
		logger.ZapSugar.Fatal(err)
	}
	defer repo.Close()

	srvGRPC := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors.HandlerWithAuth),
		grpc.ChainUnaryInterceptor(interceptors.HandlerWithLogging))
	pb.RegisterInfoKeeperServer(srvGRPC, grpcserver.NewKeeperServer(repo, *cfg))

	idleConnsClosed := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigs
		srvGRPC.GracefulStop()
		close(idleConnsClosed)
	}()

	listen, err := net.Listen("tcp", cfg.GRPC)
	if err != nil {
		logger.ZapSugar.Fatalw(err.Error(), "event", "listen port")
	}
	if err = srvGRPC.Serve(listen); err != nil {
		logger.ZapSugar.Fatalw(err.Error(), "event", "start gRPC server")
	}

	<-idleConnsClosed
}
