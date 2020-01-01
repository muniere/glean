package pubsub

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/app/server/spider"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

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
	scrape := func(uri *url.URL) error {
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

	clutch := func(uri *url.URL) error {
		uris, err := spider.Walk(uri, spider.WalkOptions{})
		if err != nil {
			return err
		}

		return spider.Download(uris, spider.DownloadOptions{
			Prefix:      config.Prefix,
			Concurrency: config.Concurrency,
			Blocking:    false,
			Overwrite:   config.Overwrite,
			DryRun:      config.DryRun,
			Interval:    500 * time.Millisecond,
		})
	}

	action := func(job task.Job, meta task.Meta) error {
		lumber.Info(box.Dict{
			"job":  job,
			"meta": meta,
		})

		uri, err := url.Parse(job.URI)
		if err != nil {
			return err
		}

		switch job.Kind {
		case rpc.Scrape:
			return scrape(uri)
		case rpc.Clutch:
			return clutch(uri)
		default:
			return errors.New(fmt.Sprintf("operation not supported: %s", job.Kind))
		}
	}

	recovery := func(err error) {
		log.Error(err)
	}

	interval := 5 * time.Second

	m.guild.Spawn(m.queue, action, recovery, interval)
}

func (m *Consumer) Start() error {
	lumber.Info(box.Dict{
		"module": "consumer",
		"action": "start",
	})
	return m.guild.Start()
}

func (m *Consumer) Stop() error {
	lumber.Info(box.Dict{
		"module": "consumer",
		"action": "stop",
	})
	return m.guild.Stop()
}

func (m *Consumer) Wait() {
	m.guild.Wait()
}
