package pubsub

import (
	"net/url"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/spider"
	"github.com/muniere/glean/internal/pkg/task"
)

type dict map[string]interface{}

type Consumer struct {
	guild *task.Guild
	queue *task.Queue
}

type ConsumerConfig struct {
	Prefix      string
	Parallel    int
	Concurrency int
	Overwrite   bool
	DryRun      bool
}

func NewConsumer(queue *task.Queue, config ConsumerConfig) *Consumer {
	s := &Consumer{
		guild: task.NewGuild(),
		queue: queue,
	}

	for i := 0; i < config.Parallel; i++ {
		s.Spawn(config)
	}

	return s
}

func (m *Consumer) Spawn(config ConsumerConfig) {
	action := func(job task.Job, meta task.Meta) error {
		log.Info(jsonic.MustEncode(
			dict{"job": job, "meta": meta},
		))

		uri, err := url.Parse(job.URI)
		if err != nil {
			return err
		}

		info, err := spider.Index(uri, spider.IndexOptions{})
		if err != nil {
			return err
		}

		return spider.Download(info.Links, spider.DownloadOptions{
			Prefix:      path.Join(config.Prefix, info.Title),
			Concurrency: config.Concurrency,
			Blocking:    false,
			Overwrite:   config.Overwrite,
			DryRun:      config.DryRun,
			Interval:    500 * time.Millisecond,
		})
	}

	recovery := func(err error) {
		log.Error(err)
	}

	interval := 5 * time.Second

	m.guild.Spawn(m.queue, action, recovery, interval)
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
