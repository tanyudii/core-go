package dwarfs

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
	Shutdown(ctx context.Context) error
	RunGracefully(t int)
	GetQueue() taskq.Queue
	RegisterWorker(workers ...Worker)
	AddMessage(ctx context.Context, name string, args ...interface{}) error
}

type service struct {
	cfg   *Config
	queue taskq.Queue
}

func NewService(factory taskq.Factory, args ...ConfigFunc) Service {
	cfg := generateConfig(args...)
	return &service{
		cfg:   cfg,
		queue: factory.RegisterQueue(cfg.ToQueueOptions()),
	}
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

func (s *service) RegisterWorker(workers ...Worker) {
	for _, worker := range workers {
		for _, task := range worker.GetTasks() {
			taskq.RegisterTask(task)
		}
	}
}

func (s *service) AddMessage(ctx context.Context, taskName string, args ...interface{}) error {
	msg := taskq.NewMessage(ctx, args...)
	msg.TaskName = taskName
	return s.queue.Add(msg)
}
