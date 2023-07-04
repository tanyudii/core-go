package zeus

import (
	"context"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	acceptlanguage "github.com/tanyudii/core-go/grpc/interceptors/accept-language"
	"github.com/tanyudii/core-go/grpc/interceptors/recovery"
	"github.com/tanyudii/core-go/grpc/interceptors/requestid"
	muxmiddleware "github.com/tanyudii/core-go/mux/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net/http"
)

func (s *service) initInterceptors() {
	s.RegisterUnaryServerInterceptor(
		requestid.UnaryServerInterceptor(),
		recovery.UnaryServerInterceptor(),
		acceptlanguage.UnaryServerInterceptor(),
	)
}

func (s *service) initConfigRestServeMuxOpts() {
	s.cfg.restServeMuxOpts = append(
		s.cfg.restServeMuxOpts,
		runtime.WithRoutingErrorHandler(muxmiddleware.MuxHandleRoutingError),
		runtime.WithErrorHandler(muxmiddleware.MuxErrorHandler),
	)
}

func (s *service) initGRPCServer() {
	s.server = grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(s.interceptors.serverUnary...)))
}

func (s *service) initReflection() {
	reflection.Register(s.GetServer())
}

func (s *service) initDefaultPrometheusCollectors() {
	s.prometheusCollectors = append(s.prometheusCollectors, RpcDurationsHistogram)
}

func (s *service) initRESTHandler(ctx context.Context) (http.Handler, error) {
	mux := runtime.NewServeMux(s.cfg.restServeMuxOpts...)

	conn, err := s.dialSelf()
	if err != nil {
		return nil, err
	}

	if err = s.initHealthCheck(mux, conn); err != nil {
		return nil, err
	}

	endpoint := ":" + s.cfg.gRPCPort
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for i := range s.restHandlers {
		h := s.restHandlers[i]
		if err := h(ctx, mux, endpoint, opts); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

func (s *service) initHealthCheck(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return mux.HandlePath(http.MethodGet, "/_health", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/plain")
		if state := conn.GetState(); state != connectivity.Ready {
			http.Error(w, fmt.Sprintf("gRPC server is %s", state), http.StatusBadGateway)
			return
		}
	})
}

func (s *service) dialSelf() (*grpc.ClientConn, error) {
	return dial("tcp", fmt.Sprintf("localhost:%s", s.cfg.gRPCPort))
}
