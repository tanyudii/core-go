package requestid

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
	"google.golang.org/grpc"
	"strings"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		md := ectx.FromIncoming(ctx)
		if acceptLang := md.Get(strings.ToLower("grpcgateway-" + ectx.RequestHeaderKeyAcceptLanguage)); acceptLang != "" {
			md.Set(strings.ToLower(ectx.RequestHeaderKeyAcceptLanguage), acceptLang)
			ctx = ectx.NewContext(ctx, ectx.NewEContext(md))
		}
		return handler(ctx, req)
	}
}
