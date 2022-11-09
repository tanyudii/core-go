package zeus

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	ginmiddleware "github.com/tanyudii/core-go/gin/middleware"
	"github.com/tanyudii/core-go/logger"
	"github.com/tanyudii/core-go/waitgroup"
	"google.golang.org/grpc"
)

type RESTHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type Service interface {
	Init()
	Shutdown(ctx context.Context) error
	GetServer() *grpc.Server
	RunGracefully(t int)
	RunServers(ctx context.Context) <-chan error
	ListenAndServeGRPC(ctx context.Context) error
	ListenAndServeREST(ctx context.Context) error
	RegisterUnaryServerInterceptor(i ...grpc.UnaryServerInterceptor)
	RegisterRESTHandler(handlers ...RESTHandler)
}

type service struct {
	cfg          *Config
	server       *grpc.Server
	interceptors Interceptors
	restHandlers []RESTHandler
}

type Interceptors struct {
	serverUnary []grpc.UnaryServerInterceptor
}

func NewService(args ...ConfigFunc) Service {
	return &service{
		cfg: generateConfig(args...),
	}
}

func (s *service) Init() {
	s.initInterceptors()
	s.initConfigRestServeMuxOpts()
	s.initGRPCServer()
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) RunGracefully(t int) {
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	go func() {
		if err := <-s.RunServers(mainCtx); err != nil {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Infof("Server is shutting down: for %ds %v", t, time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()
	cancelMainCtx()
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown err: %v", err)
	}
	logger.Infof("Server exiting %v", time.Now())
}

func (s *service) RunServers(ctx context.Context) <-chan error {
	var once sync.Once
	exitCh := make(chan error)
	waitGroup := &waitgroup.WgWrapper{}
	exitFunc := func(err error) {
		once.Do(func() {
			exitCh <- err
		})
	}

	go waitGroup.Wrap(func() {
		logger.Infof("Initializing gRPC connection in port %s", s.cfg.gRPCPort)
		exitFunc(s.ListenAndServeGRPC(ctx))
	})

	go waitGroup.Wrap(func() {
		logger.Infof("Initializing HTTP connection in port %s", s.cfg.restPort)
		exitFunc(s.ListenAndServeREST(ctx))
	})

	return exitCh
}

func (s *service) ListenAndServeGRPC(_ context.Context) error {
	if s.server == nil {
		return errors.New("ListenAndServeGRPC: server is not initialized")
	}
	logger.Infof("starting gRPC server at :%s...", s.cfg.gRPCPort)

	defer s.server.GracefulStop()
	lis, err := net.Listen("tcp", ":"+s.cfg.gRPCPort)
	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func (s *service) ListenAndServeREST(ctx context.Context) error {
	handler, err := s.initRESTHandler(ctx)
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	srv := &http.Server{
		Addr:    ":" + s.cfg.restPort,
		Handler: r,
	}

	//register CORS when config enabled
	if s.cfg.enableCORS {
		r.Use(ginmiddleware.GinCORS())
	}

	//register JSON when config enabled
	if s.cfg.onlyJSON {
		r.Use(ginmiddleware.GinJSON())
	}

	r.Group("*{any}").Any("", gin.WrapH(handler))

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			logger.Errorf("ListenAndServeREST: failed to shutdown %v\n", err)
		}
	}()

	logger.Infof("starting HTTP server at :%s...", s.cfg.restPort)
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("ListenAndServeREST: failed to listen and serve: %v\n", err)
		return err
	}

	return nil
}

func (s *service) GetServer() *grpc.Server {
	return s.server
}

func (s *service) RegisterUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) {
	s.interceptors.serverUnary = append(s.interceptors.serverUnary, interceptors...)
}

func (s *service) RegisterRESTHandler(handlers ...RESTHandler) {
	s.restHandlers = append(s.restHandlers, handlers...)
}
