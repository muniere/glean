package pubsub

import (
	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Config struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	Prefix      string `json:"prefix"`
	Parallel    int    `json:"parallel"`
	Concurrency int    `json:"concurrency"`
	Overwrite   bool   `json:"overwrite"`
	LogDir      string `json:"log_dir"`
	DryRun      bool   `json:"dry_run"`
	Verbose     bool   `json:"verbose"`
}

func NewSupervisor(config Config) *Supervisor {
	queue := task.NewQueue()

	consumer := NewConsumer(queue, ConsumerConfig{
		Parallel:    config.Parallel,
		Concurrency: config.Concurrency,
		Prefix:      config.Prefix,
		Overwrite:   config.Overwrite,
		DryRun:      config.DryRun,
	})

	producer := NewProducer(queue, ProducerConfig{
		Address: config.Address,
		Port:    config.Port,
	})

	producer.Register(rpc.Config, func(c *action.Context) error {
		return c.Gateway.Success(config)
	})

	return &Supervisor{
		queue:    queue,
		producer: producer,
		consumer: consumer,
	}
}

type Supervisor struct {
	queue    *task.Queue
	producer *Producer
	consumer *Consumer
}

func (s *Supervisor) Start() error {
	if err := s.consumer.Start(); err != nil {
		return err
	}
	if err := s.producer.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Supervisor) Stop() error {
	if err := s.consumer.Stop(); err != nil {
		return err
	}
	if err := s.producer.Stop(); err != nil {
		return err
	}
	return nil
}
