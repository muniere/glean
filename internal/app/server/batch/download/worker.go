package download

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/muniere/glean/internal/app/server/batch/lumber"
	"github.com/muniere/glean/internal/pkg/images"
	"github.com/muniere/glean/internal/pkg/path"
	"github.com/muniere/glean/internal/pkg/sys"
)

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
			if err == images.TooSmall {
				continue
			}
			if err == images.TooLarge {
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

	if err := precond(ctx, options); err != nil {
		if err == skipDownload {
			lumber.SkipStep("fetch", ctx.dict())
		}
		return err
	}

	temp, err := fetch(ctx, options)
	if err != nil {
		return err
	}

	ctx.temp = temp.String()
	if err := postcond(ctx, options); err != nil {
		if err == images.TooSmall {
			lumber.SkipStep("move", ctx.dict())
		}
		return err
	}

	return move(ctx, options)
}

func compose(cmd command, options Options) context {
	p := path.New(cmd.uri.String()).Base()

	if len(options.Prefix) > 0 {
		p = p.Prepend(options.Prefix)
	}

	return context{
		uri:  cmd.uri,
		temp: "",
		path: p.String(),
	}
}

func precond(ctx context, options Options) error {
	if options.DryRun {
		return skipDownload
	}
	if !options.Overwrite && sys.Exists(ctx.path) {
		return skipDownload
	}
	return nil
}

func fetch(ctx context, options Options) (*path.Pathname, error) {
	lumber.Start(ctx.dict())
	defer lumber.Finish(ctx.dict())
	return doFetch(ctx.uri, path.Base(ctx.path))
}

func postcond(ctx context, options Options) error {
	return images.TestFile(ctx.temp, options.Scope)
}

func move(ctx context, options Options) error {
	lumber.Start(ctx.dict())
	defer lumber.Finish(ctx.dict())
	return doMove(ctx.temp, ctx.path)
}

//
// Helper
//
func doFetch(uri *url.URL, name string) (*path.Pathname, error) {
	res, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	f, err := ioutil.TempFile("", name)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return nil, err
	}

	return path.New(f.Name()), err
}

func doMove(src string, dest string) error {
	return os.Rename(src, dest)
}
