package zeus

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultMaxCallRcvMsgSize  = 1024 * 1024 * 30
	defaultMaxCallSendMsgSize = 1024 * 1024 * 30
)

func ClientConn(addr string) grpc.ClientConnInterface {
	opts := getDialOpts()
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	return conn
}

func getDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultMaxCallRcvMsgSize),
			grpc.MaxCallSendMsgSize(defaultMaxCallSendMsgSize),
		),
	}
}

func dial(network, addr string) (*grpc.ClientConn, error) {
	return dialContext(context.Background(), network, addr)
}

func dialContext(ctx context.Context, network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return dialTCP(ctx, addr)
	default:
		return nil, fmt.Errorf("unsupported network type %q", network)
	}
}

func dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	opts := getDialOpts()
	return grpc.DialContext(ctx, addr, opts...)
}
