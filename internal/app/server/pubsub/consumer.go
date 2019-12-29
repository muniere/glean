package pubsub

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/task"
)

type Consumer struct {
	guild *task.Guild
	queue *task.Queue
}

type ConsumerConfig struct {
	Concurrency int
}

func NewConsumer(queue *task.Queue, config ConsumerConfig) *Consumer {
	s := &Consumer{
		guild: task.NewGuild(),
		queue: queue,
	}

	for i := 0; i < config.Concurrency; i++ {
		s.Spawn()
	}

	return s
}

func (m *Consumer) Spawn() {
	m.guild.Spawn(
		m.queue,
		func(job task.Job, meta task.Meta) error {
			// TODO: Crawl with query
			m := map[string]interface{}{
				"job":  job,
				"meta": meta,
			}
			b, _ := json.Marshal(m)
			log.Info(string(b))
			return nil
		},
		func(err error) {
			log.Error(err)
		},
		5*time.Second,
	)
}

func (m *Consumer) Start() error {
	return m.guild.Start()
}

func (m *Consumer) Stop() error {
	return m.guild.Stop()
}

func (m *Consumer) Wait() {
	m.guild.Wait()
}
