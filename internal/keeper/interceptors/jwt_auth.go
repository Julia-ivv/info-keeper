package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Julia-ivv/info-keeper.git/internal/authorizer"
	pb "github.com/Julia-ivv/info-keeper.git/internal/proto/pb"
)

// HandlerWithAuth adds user token to the handler.
func HandlerWithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == pb.InfoKeeper_AddUser_FullMethodName ||
		info.FullMethod == pb.InfoKeeper_AuthUser_FullMethodName {
		return handler(ctx, req)
	}

	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(authorizer.AccessToken)
		if len(values) > 0 {
			token = values[0]
		}
	}
	if len(token) == 0 {
		return nil, status.Error(codes.Internal, "missing token")
	}

	ctx = context.WithValue(ctx, authorizer.UserContextKey, token)

	return handler(ctx, req)
}
