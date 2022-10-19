package recovery

import (
	"context"
	"runtime/debug"

	"github.com/tanyudii/core-go/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Unknown, "unexpected error happened")
				logger.WithField("stacktrace", string(debug.Stack())).Errorf("panic recovered: %v", r)
			}
		}()
		return handler(ctx, req)
	}
}
