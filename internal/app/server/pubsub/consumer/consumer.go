package consumer

import (
	"time"

	"github.com/muniere/glean/internal/pkg/lumber"
	. "github.com/muniere/glean/internal/pkg/stdlib"
	"github.com/muniere/glean/internal/pkg/task"
)

type Consumer struct {
	guild *task.Guild
	queue *task.Queue
}

type Config struct {
	DataDir     string
	Parallel    int
	Concurrency int
	MinWidth    int
	MaxWidth    int
	MinHeight   int
	MaxHeight   int
	Overwrite   bool
	DryRun      bool
}

func NewConsumer(queue *task.Queue, config Config) *Consumer {
	x := &Consumer{
		guild: task.NewGuild(),
		queue: queue,
	}

	for i := 0; i < config.Parallel; i++ {
		x.Spawn(config)
	}

	return x
}

func (x *Consumer) Spawn(config Config) {
	x.guild.Spawn(
		x.queue,
		&actionAdapter{config},
		&recoveryAdapter{config},
		5*time.Second,
	)
}

func (x *Consumer) Start() error {
	lumber.Info(NewDict(Pair("module", "consumer"), Pair("event", "start")))
	return x.guild.Start()
}

func (x *Consumer) Stop() error {
	lumber.Info(NewDict(Pair("module", "consumer"), Pair("event", "stop")))
	return x.guild.Stop()
}

func (x *Consumer) Wait() {
	x.guild.Wait()
}
