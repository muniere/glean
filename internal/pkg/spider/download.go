package spider

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/sys"
)

var SkipDownload = errors.New("skip run")

type DownloadOptions struct {
	Prefix      string
	Concurrency int
	Blocking    bool
	Overwrite   bool
	DryRun      bool
	Interval    time.Duration
}

type command struct {
	uri *url.URL
}

func Download(urls []*url.URL, options DownloadOptions) error {
	// prepare
	if !options.DryRun && !sys.Exists(options.Prefix) {
		log.Debug(jsonic.MustEncode(dict{
			"label":      "start",
			"action":     "mkdir",
			"param.path": options.Prefix,
		}))
		if err := os.MkdirAll(options.Prefix, 0755); err != nil {
			return err
		}
	}

	// workers
	wg := &sync.WaitGroup{}
	ch := make(chan command, len(urls))

	for i := 0; i < options.Concurrency; i++ {
		launch(wg, ch, options)
	}

	// enqueue
	for _, u := range urls {
		ch <- command{uri: u}
	}
	close(ch)

	// join
	wg.Wait()

	return nil
}

func launch(group *sync.WaitGroup, channel chan command, options DownloadOptions) {
	group.Add(1)

	go func() {
		defer group.Done()

		for {
			cmd, ok := <-channel
			if !ok {
				return
			}

			err := run(cmd, options)
			if err == SkipDownload {
				continue
			}
			if err != nil {
				log.Warn(err)
			}

			time.Sleep(options.Interval)
		}
	}()
}

func run(cmd command, options DownloadOptions) error {
	log.Debug(jsonic.MustEncode(dict{
		"label":     "start",
		"action":    "run",
		"param.uri": cmd.uri.String(),
	}))

	base := filepath.Base(cmd.uri.String())

	var path string

	if options.Prefix != "" {
		path = filepath.Join(options.Prefix, base)
	} else {
		path = base
	}

	if !options.Overwrite && sys.Exists(path) {
		log.Info(jsonic.MustEncode(dict{
			"label":      "skip",
			"action":     "run",
			"param.uri":  cmd.uri.String(),
			"param.path": path,
		}))
		return SkipDownload
	}

	log.Info(jsonic.MustEncode(dict{
		"label":      "start",
		"action":     "download",
		"param.uri":  cmd.uri.String(),
		"param.path": path,
	}))

	if options.DryRun {
		log.Info(jsonic.MustEncode(dict{
			"label":      "skip",
			"action":     "download",
			"param.uri":  cmd.uri.String(),
			"param.path": path,
		}))
		return SkipDownload
	}

	log.Debug(jsonic.MustEncode(dict{
		"label":     "start",
		"action":    "fetch",
		"param.uri": cmd.uri.String(),
	}))

	res, err := http.Get(cmd.uri.String())
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	log.Debug(jsonic.MustEncode(dict{
		"label":      "start",
		"action":     "create",
		"param.path": path,
	}))

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	log.Debug(jsonic.MustEncode(dict{
		"label":     "start",
		"action":    "copy",
		"param.src": cmd.uri.String(),
		"param.dst": path,
	}))

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	log.Info(jsonic.MustEncode(dict{
		"label":     "finish",
		"action":    "run",
		"param.src": cmd.uri.String(),
		"param.dst": path,
	}))

	return nil
}
