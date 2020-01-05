package download

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/muniere/glean/internal/app/server/batch/lumber"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/sys"
)

var skipDownload = errors.New("skip download")

type command struct {
	uri *url.URL
}

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
	ctx := box.Dict{"path": options.Prefix}

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

func launch(group *sync.WaitGroup, channel chan command, options Options) {
	group.Add(1)

	go func() {
		defer group.Done()

		for {
			cmd, ok := <-channel
			if !ok {
				return
			}

			err := run(cmd, options)
			if err == skipDownload {
				continue
			}
			if err != nil {
				lumber.Warn(err)
			}

			time.Sleep(options.Interval)
		}
	}()
}

func run(cmd command, options Options) error {
	ctx := compose(cmd, options)

	if err := test(ctx, options); err == skipDownload {
		return skipDownload
	}

	res, err := fetch(ctx, options)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	file, err := touch(ctx, options)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	return save(file, res.Body, ctx, options)
}

func compose(cmd command, options Options) context {
	base := filepath.Base(cmd.uri.String())

	var path string

	if options.Prefix != "" {
		path = filepath.Join(options.Prefix, base)
	} else {
		path = base
	}

	return context{uri: cmd.uri, path: path}
}

func test(context context, options Options) error {
	if options.DryRun {
		lumber.SkipStep("download", context.dict())
		return skipDownload
	}

	if !options.Overwrite && sys.Exists(context.path) {
		lumber.SkipStep("download", context.dict())
		return skipDownload
	}

	return nil
}

func fetch(context context, options Options) (*http.Response, error) {
	lumber.Start(context.dict())

	defer lumber.Finish(context.dict())

	return http.Get(context.uri.String())
}

func touch(context context, options Options) (*os.File, error) {
	lumber.Start(context.dict())

	defer lumber.Finish(context.dict())

	return os.Create(context.path)
}

func save(dst io.Writer, src io.Reader, context context, options Options) error {
	lumber.Start(context.dict())

	defer lumber.Start(context.dict())

	_, err := io.Copy(dst, src)
	return err
}
