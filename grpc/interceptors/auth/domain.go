package auth

import (
	"context"
	"google.golang.org/grpc"
)

type Service interface {
	authenticate(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error)
	authenticateToken(ctx context.Context) (context.Context, error)
	authorizedInternalCall(ctx context.Context) (context.Context, bool)
}
