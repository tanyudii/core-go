package auth

import (
	"context"
	"google.golang.org/grpc"
)

func UnaryInterceptor(
	tokenService TokenService,
	args ...ConfigFunc,
) grpc.UnaryServerInterceptor {
	svc := newService(tokenService, args...)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := svc.authenticate(ctx, info)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}
