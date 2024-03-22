package main

import (
	"google.golang.org/grpc"

	"github.com/Julia-ivv/info-keeper.git/internal/config"
	"github.com/Julia-ivv/info-keeper.git/internal/grpcserver"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto"
	"github.com/Julia-ivv/info-keeper.git/internal/storage"
	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

func main() {
	cfg := config.NewConfig()

	logger.ZapSugar = logger.NewLogger()
	logger.ZapSugar.Infow("Starting server", "addr", cfg.Host)
	logger.ZapSugar.Infow("flags", "db dsn", cfg.DBDSN)

	repo, err := storage.NewStorage(*cfg)
	if err != nil {
		logger.ZapSugar.Fatal(err)
	}
	defer repo.Close()

	// var srv = http.Server{
	// 	Addr: cfg.Host,
	// 	//	Handler: httpserver.NewURLRouter(repo, *cfg, &httpWg),
	// }

	// certFile, privateKeyFile, err := certgenerator.GenCert(4096)
	// if err != nil {
	// 	logger.ZapSugar.Fatalw(err.Error(), "event", "create certificate or private key")
	// }
	// err = srv.ListenAndServeTLS(certFile.Name(), privateKeyFile.Name())
	// if err != nil && err != http.ErrServerClosed {
	// 	logger.ZapSugar.Fatalw(err.Error(), "event", "start server")
	// }

	srvGRPC := grpc.NewServer()
	// grpc.ChainUnaryInterceptor(interceptors.HandlerWithAuth),
	// grpc.ChainUnaryInterceptor(interceptors.HandlerWithLogging))
	pb.RegisterInfoKeeperServer(srvGRPC, grpcserver.NewShortenerServer(repo, *cfg))
}
