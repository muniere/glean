package consumer

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/muniere/glean/internal/app/server/batch"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
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
	scrape := func(uri *url.URL, prefix string) error {
		info, err := batch.Index(uri, batch.IndexOptions{})
		if err != nil {
			return err
		}

		var prefixer string
		if len(prefix) > 0 {
			if filepath.IsAbs(prefix) {
				prefixer = prefix
			} else {
				prefixer = filepath.Join(config.DataDir, prefix)
			}
		} else {
			if len(info.Title) > 0 {
				prefixer = filepath.Join(config.DataDir, info.Title)
			} else {
				prefixer = filepath.Join(config.DataDir, url.QueryEscape(uri.String()))
			}
		}

		opts := batch.DownloadOptions{
			Prefix:      prefixer,
			Concurrency: config.Concurrency,
			MinWidth:    config.MinWidth,
			MaxWidth:    config.MaxWidth,
			MinHeight:   config.MinHeight,
			MaxHeight:   config.MaxHeight,
			Blocking:    false,
			Overwrite:   config.Overwrite,
			DryRun:      config.DryRun,
			Interval:    500 * time.Millisecond,
		}

		return batch.Download(info.Links, opts)
	}

	clutch := func(uri *url.URL, prefix string) error {
		uris, err := batch.Walk(uri, batch.WalkOptions{})
		if err != nil {
			return err
		}

		var prefixer string
		if len(prefix) > 0 {
			if filepath.IsAbs(prefix) {
				prefixer = prefix
			} else {
				prefixer = filepath.Join(config.DataDir, prefix)
			}
		} else {
			prefixer = filepath.Join(config.DataDir, url.QueryEscape(uri.String()))
		}

		opts := batch.DownloadOptions{
			Prefix:      prefixer,
			Concurrency: config.Concurrency,
			MinWidth:    config.MinWidth,
			MaxWidth:    config.MaxWidth,
			MinHeight:   config.MinHeight,
			MaxHeight:   config.MaxHeight,
			Blocking:    false,
			Overwrite:   config.Overwrite,
			DryRun:      config.DryRun,
			Interval:    500 * time.Millisecond,
		}

		return batch.Download(uris, opts)
	}

	action := func(job task.Job, meta task.Meta) error {
		lumber.Info(box.Dict{
			"module": "consumer",
			"event":  "job::consume",
			"job":    job,
			"meta":   meta,
		})

		uri, err := url.Parse(job.URI)
		if err != nil {
			return err
		}

		switch job.Kind {
		case rpc.Scrape:
			return scrape(uri, job.Prefix)
		case rpc.Clutch:
			return clutch(uri, job.Prefix)
		default:
			return errors.New(fmt.Sprintf("operation not supported: %s", job.Kind))
		}
	}

	recovery := func(err error) {
		lumber.Error(box.Dict{
			"module": "consumer",
			"event":  "error",
			"error":  err.Error(),
		})
	}

	interval := 5 * time.Second

	x.guild.Spawn(x.queue, action, recovery, interval)
}

func (x *Consumer) Start() error {
	lumber.Info(box.Dict{
		"module": "consumer",
		"event":  "start",
	})
	return x.guild.Start()
}

func (x *Consumer) Stop() error {
	lumber.Info(box.Dict{
		"module": "consumer",
		"event":  "stop",
	})
	return x.guild.Stop()
}

func (x *Consumer) Wait() {
	x.guild.Wait()
}
