package grpcserver

import (
	"github.com/Julia-ivv/info-keeper.git/internal/config"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto"
	"github.com/Julia-ivv/info-keeper.git/internal/storage"
)

// ShortenerServer stores the repository and settings of this application.
type KeeperGRPCServer struct {
	pb.UnimplementedInfoKeeperServer
	stor storage.Repositorier
	cfg  config.Flags
}

// NewShortenerServer creates an instance with storage and settings for grpc methods.
func NewShortenerServer(stor storage.Repositorier, cfg config.Flags) *KeeperGRPCServer {
	k := &KeeperGRPCServer{}
	k.stor = stor
	k.cfg = cfg
	return k
}
