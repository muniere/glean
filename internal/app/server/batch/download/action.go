package download

import (
	"errors"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/muniere/glean/internal/app/server/batch/lumber"
	"github.com/muniere/glean/internal/pkg/images"
	"github.com/muniere/glean/internal/pkg/std"
	"github.com/muniere/glean/internal/pkg/sys"
)

//
// Error
//
var skipDownload = errors.New("skip download")

//
// Struct
//
type Options struct {
	Prefix      string
	Scope       images.Scope
	Interval    time.Duration
	Concurrency int
	Blocking    bool
	Overwrite   bool
	DryRun      bool
}

type command struct {
	uri *url.URL
}

type context struct {
	uri  *url.URL
	temp string
	path string
}

func (c *context) dict() std.Dict {
	dict := std.NewDict()
	if c.uri != nil {
		dict.Put("uri", c.uri.String())
	}
	if len(c.temp) > 0 {
		dict.Put("temp", c.temp)
	}
	if len(c.path) > 0 {
		dict.Put("path", c.path)
	}
	return dict
}

//
// Action
//
func Perform(urls []*url.URL, options Options) error {
	// prepare
	if err := mkdir(options); err != nil {
		return err
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

func mkdir(options Options) error {
	ctx := std.NewDict(
		std.Pair("path", options.Prefix),
	)

	if options.DryRun {
		lumber.Skip(ctx)
		return nil
	}

	if sys.Exists(options.Prefix) {
		lumber.Skip(ctx)
		return nil
	}

	lumber.Start(ctx)

	defer lumber.Finish(ctx)

	return os.MkdirAll(options.Prefix, 0755)
}
