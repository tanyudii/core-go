package auth

import (
	"context"
	"github.com/tanyudii/core-go/auth"
	"google.golang.org/grpc"
)

func UnaryInterceptor(
	authService auth.Service,
	tokenService auth.TokenService,
	args ...ConfigFunc,
) grpc.UnaryServerInterceptor {
	svc := newService(authService, tokenService, args...)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := svc.authenticate(ctx, info)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}
