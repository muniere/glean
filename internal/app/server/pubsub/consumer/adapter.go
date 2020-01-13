package consumer

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/muniere/glean/internal/app/server/batch"
	"github.com/muniere/glean/internal/pkg/images"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/path"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/std"
	"github.com/muniere/glean/internal/pkg/task"
)

//
// Adapter for task.Action
//
type actionAdapter struct {
	config Config
}

func (x *actionAdapter) Invoke(job task.Job, meta task.Meta) error {
	lumber.Info(std.NewDict(std.Pair("module", "consumer"), std.Pair("event", "job::consume"), std.Pair("job", job), std.Pair("meta", meta)))

	uri, err := url.Parse(job.URI)
	if err != nil {
		return err
	}

	switch job.Kind {
	case rpc.Scrape:
		return x.scrape(uri, job.Prefix)
	case rpc.Clutch:
		return x.clutch(uri, job.Prefix)
	default:
		return errors.New(fmt.Sprintf("operation not supported: %s", job.Kind))
	}
}

func (x *actionAdapter) scrape(uri *url.URL, prefix string) error {
	info, err := batch.Index(uri, batch.IndexOptions{})
	if err != nil {
		return err
	}

	pref := func() *path.Pathname {
		if len(prefix) > 0 && path.IsAbs(prefix) {
			return path.New(prefix)
		}
		if len(prefix) > 0 {
			return path.New(x.config.DataDir).Append(prefix)
		}
		if len(info.Title) > 0 {
			return path.New(x.config.DataDir).Append(info.Title)
		} else {
			return path.New(x.config.DataDir).Append(url.QueryEscape(uri.String()))
		}
	}()

	scope := images.Scope{
		Min: images.Size{
			Width:  x.config.MinWidth,
			Height: x.config.MinHeight,
		},
		Max: images.Size{
			Width:  x.config.MaxWidth,
			Height: x.config.MaxHeight,
		},
	}

	opts := batch.DownloadOptions{
		Prefix:      pref.String(),
		Scope:       scope,
		Interval:    500 * time.Millisecond,
		Concurrency: x.config.Concurrency,
		Blocking:    false,
		Overwrite:   x.config.Overwrite,
		DryRun:      x.config.DryRun,
	}

	return batch.Download(info.Links, opts)
}

func (x *actionAdapter) clutch(uri *url.URL, prefix string) error {
	uris, err := batch.Walk(uri, batch.WalkOptions{})
	if err != nil {
		return err
	}

	pref := func() *path.Pathname {
		if len(prefix) > 0 && path.IsAbs(prefix) {
			return path.New(prefix)
		}
		if len(prefix) > 0 {
			return path.New(x.config.DataDir).Append(prefix)
		}
		return path.New(x.config.DataDir).Append(url.QueryEscape(uri.String()))
	}()

	scope := images.Scope{
		Min: images.Size{
			Width:  x.config.MinWidth,
			Height: x.config.MinHeight,
		},
		Max: images.Size{
			Width:  x.config.MaxWidth,
			Height: x.config.MaxHeight,
		},
	}

	opts := batch.DownloadOptions{
		Prefix:      pref.String(),
		Scope:       scope,
		Interval:    500 * time.Millisecond,
		Concurrency: x.config.Concurrency,
		Blocking:    false,
		Overwrite:   x.config.Overwrite,
		DryRun:      x.config.DryRun,
	}

	return batch.Download(uris, opts)
}

//
// Adapter for task.Recovery
//
type recoveryAdapter struct {
	config Config
}

func (x *recoveryAdapter) Invoke(err error) {
	lumber.Error(std.NewDict(std.Pair("module", "consumer"), std.Pair("event", "error"), std.Pair("error", err.Error())))
}
