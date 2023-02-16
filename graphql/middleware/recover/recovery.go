package recover

import (
	"context"
	"github.com/tanyudii/core-go/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func Recover(ctx context.Context, err interface{}) error {
	err = status.Errorf(codes.Unknown, "unexpected error happened")
	logger.WithField("stacktrace", string(debug.Stack())).Errorf("panic recovered: %v", err)
	return gqlerror.Errorf("internal system error")
}
