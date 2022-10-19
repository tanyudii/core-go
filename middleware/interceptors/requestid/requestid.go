package requestid

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	ectx "github.com/tanyudii/core-go/econtext"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		md := ectx.FromIncoming(ctx)
		if md.Get("requestid") == "" {
			requestID := fmt.Sprintf("%s-%d", uuid.NewString(), time.Now().Unix())
			md.Set("requestid", requestID)
			ctx = md.ToIncoming(ctx)
		}
		return handler(ctx, req)
	}
}
