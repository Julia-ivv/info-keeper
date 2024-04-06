package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

// HandlerWithLogging adds logging to the gRPC methods.
func HandlerWithLogging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	logger.ZapSugar.Infoln(
		"full method", info.FullMethod,
		"duration", time.Since(start),
	)

	return h, err
}
