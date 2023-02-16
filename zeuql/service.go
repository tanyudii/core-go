package zeuql

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginmiddleware "github.com/tanyudii/core-go/gin/middleware"
	"github.com/tanyudii/core-go/logger"
	muxmiddleware "github.com/tanyudii/core-go/mux/middleware"
	"github.com/tanyudii/core-go/waitgroup"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Service interface {
	Init()
	Shutdown(ctx context.Context) error
	GetEngine() *gin.Engine
	RunGracefully(t int)
	RunServers(ctx context.Context) <-chan error
	RegisterExecutableSchema(schema graphql.ExecutableSchema)
}

type service struct {
	cfg                  *Config
	engine               *gin.Engine
	schema               graphql.ExecutableSchema
	prometheusCollectors []prometheus.Collector
}

func NewService(args ...ConfigFunc) Service {
	return &service{
		cfg: generateConfig(args...),
	}
}

func (s *service) Init() {
	s.initEngine()
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) GetEngine() *gin.Engine {
	return s.engine
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
		logger.Infof("Initializing graphQL connection in port %s", s.cfg.graphQLPort)
		exitFunc(s.ListenAndServeGraphQL(ctx))
	})

	go waitGroup.Wrap(func() {
		logger.Infof("Initializing Prometheus connection in port %s", s.cfg.prometheusPort)
		exitFunc(s.ListenAndServePrometheus(ctx))
	})

	return exitCh
}

func (s *service) ListenAndServeGraphQL(ctx context.Context) (err error) {
	if s.engine == nil {
		return errors.New("ListenAndServeGraphQL: engine is not initialized")
	}

	srv := &http.Server{
		Addr:    ":" + s.cfg.graphQLPort,
		Handler: s.engine,
	}

	//register CORS when config enabled
	if s.cfg.enableCORS {
		s.engine.Use(ginmiddleware.GinCORS())
	}

	s.initHealthCheck()

	s.engine.POST("/query", s.graphQLHandler())
	s.engine.GET("/", s.playgroundHandler())

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			logger.Errorf("ListenAndServeGraphQL: failed to shutdown %v\n", err)
		}
	}()

	logger.Infof("starting graphQL server at :%s...", s.cfg.graphQLPort)
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("ListenAndServeGraphQL: failed to listen and serve: %v\n", err)
		return err
	}

	return nil
}

func (s *service) ListenAndServePrometheus(ctx context.Context) (err error) {
	for _, c := range s.prometheusCollectors {
		if err = prometheus.Register(c); err != nil {
			logger.Errorf("ListenAndServePrometheus: failed to register collector: %v\n", err)
		}
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.prometheusPort),
		Handler: muxmiddleware.MuxCORS(mux),
	}

	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			logger.Errorf("ListenAndServePrometheus: failed to shutdown %v\n", err)
		}
	}()

	logger.Infof("starting Prometheus server at :%s...", s.cfg.prometheusPort)
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("ListenAndServePrometheus: failed to listen and serve: %v\n", err)
		return err
	}

	return nil
}

func (s *service) RegisterExecutableSchema(schema graphql.ExecutableSchema) {
	s.schema = schema
}
