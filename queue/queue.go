package queue

import (
	"context"
	"github.com/tanyudii/core-go/logger"
	"github.com/vmihailenco/taskq/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service interface {
}

type service struct {
	cfg     *Config
	factory taskq.Factory
	queue   taskq.Queue
	tasks   []taskq.TaskOptions
}

func NewService(factory taskq.Factory, args ...ConfigFunc) Service {
	return &service{
		factory: factory,
		cfg:     generateConfig(args...),
	}
}

func (s *service) Init() {
	s.initQueue()
}

func (s *service) initQueue() {
	s.queue = s.factory.RegisterQueue(s.cfg.ToQueueOptions())
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) RunGracefully(t int) {
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	go func() {
		if err := s.queue.Consumer().Start(mainCtx); err != nil {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Infof("Queue is shutting down: for %ds %v", t, time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()
	cancelMainCtx()
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatalf("Queue shutdown err: %v", err)
	}
	logger.Infof("Queue exiting %v", time.Now())
}

func (s *service) GetQueue() taskq.Queue {
	return s.queue
}

func (s *service) AddTask(tasks ...taskq.TaskOptions) {
	for i := range tasks {
		taskq.RegisterTask(&tasks[i])
	}
}
