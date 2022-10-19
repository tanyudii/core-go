package zeus

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultMaxCallRcvMsgSize  = 1024 * 1024 * 30
	defaultMaxCallSendMsgSize = 1024 * 1024 * 30
)

func ClientConn(addr string) grpc.ClientConnInterface {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultMaxCallRcvMsgSize),
			grpc.MaxCallSendMsgSize(defaultMaxCallSendMsgSize),
		),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	return conn
}
