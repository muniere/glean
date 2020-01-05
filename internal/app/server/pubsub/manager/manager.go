package manager

import (
	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/app/server/pubsub/consumer"
	"github.com/muniere/glean/internal/app/server/pubsub/producer"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Manager struct {
	queue    *task.Queue
	producer *producer.Producer
	consumer *consumer.Consumer
}

type Config struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	DataDir     string `json:"data_dir"`
	Parallel    int    `json:"parallel"`
	Concurrency int    `json:"concurrency"`
	MinWidth    int    `json:"min_width"`
	MaxWidth    int    `json:"max_width"`
	MinHeight   int    `json:"min_height"`
	MaxHeight   int    `json:"max_height"`
	Overwrite   bool   `json:"overwrite"`
	LogDir      string `json:"log_dir"`
	DryRun      bool   `json:"dry_run"`
	Verbose     bool   `json:"verbose"`
}

func NewManager(config Config) *Manager {
	q := task.NewQueue()

	c := consumer.NewConsumer(q, consumer.Config{
		Parallel:    config.Parallel,
		Concurrency: config.Concurrency,
		DataDir:     config.DataDir,
		MinWidth:    config.MinWidth,
		MaxWidth:    config.MaxWidth,
		MinHeight:   config.MinHeight,
		MaxHeight:   config.MaxHeight,
		Overwrite:   config.Overwrite,
		DryRun:      config.DryRun,
	})

	p := producer.NewProducer(q, producer.Config{
		Address: config.Address,
		Port:    config.Port,
	})

	p.Register(rpc.Config, func(c *action.Context) error {
		return c.Gateway.Success(config)
	})

	return &Manager{
		queue:    q,
		producer: p,
		consumer: c,
	}
}

func (x *Manager) Start() error {
	if err := x.consumer.Start(); err != nil {
		return err
	}
	if err := x.producer.Start(); err != nil {
		return err
	}
	return nil
}

func (x *Manager) Stop() error {
	if err := x.consumer.Stop(); err != nil {
		return err
	}
	if err := x.producer.Stop(); err != nil {
		return err
	}
	return nil
}
