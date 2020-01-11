package download

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
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

//
// Error
//
var skipDownload = errors.New("skip download")
var dataTooSmall = errors.New("data too small")
var dataTooLarge = errors.New("data too large")

//
// Struct
//
type Options struct {
	Prefix      string
	Concurrency int
	MinWidth    int
	MaxWidth    int
	MinHeight   int
	MaxHeight   int
	Blocking    bool
	Overwrite   bool
	DryRun      bool
	Interval    time.Duration
}

type command struct {
	uri *url.URL
}

type context struct {
	uri  *url.URL
	temp string
	path string
}

func (c *context) dict() box.Dict {
	dict := box.Dict{}
	if c.uri != nil {
		dict["uri"] = c.uri.String()
	}
	if len(c.temp) > 0 {
		dict["temp"] = c.temp
	}
	if len(c.path) > 0 {
		dict["path"] = c.path
	}
	return dict
}

//
// Action / Supervisor
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

//
// Action / Worker
//
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
			if err == dataTooSmall {
				continue
			}
			if err == dataTooLarge {
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

	if err := preCondition(ctx, options); err != nil {
		if err == skipDownload {
			lumber.SkipStep("download", ctx.dict())
		}
		return err
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

	temp, err := save(res.Body, ctx, options)
	if err != nil {
		return err
	}

	ctx.temp = temp

	if err := postCondition(temp, ctx, options); err != nil {
		if err == dataTooSmall {
			lumber.SkipStep("persistent", ctx.dict())
		}
		return err
	}

	return persistent(temp, ctx, options)
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

func preCondition(ctx context, options Options) error {
	if options.DryRun {
		return skipDownload
	}

	if !options.Overwrite && sys.Exists(ctx.path) {
		return skipDownload
	}

	return nil
}

func fetch(ctx context, options Options) (*http.Response, error) {
	lumber.Start(ctx.dict())

	defer lumber.Finish(ctx.dict())

	return http.Get(ctx.uri.String())
}

func save(src io.Reader, ctx context, options Options) (string, error) {
	lumber.Start(ctx.dict())

	defer lumber.Finish(ctx.dict())

	f, err := ioutil.TempFile("", filepath.Base(ctx.path))
	if err != nil {
		return "", err
	}

	defer func() {
		_ = f.Close()
	}()

	ctx.temp = f.Name()

	if err := temp(f, src, ctx, options); err != nil {
		return "", err
	}

	return f.Name(), err
}

func temp(f *os.File, src io.Reader, ctx context, options Options) error {
	lumber.Start(ctx.dict())

	defer lumber.Finish(ctx.dict())

	_, err := io.Copy(f, src)
	return err
}

func postCondition(path string, ctx context, options Options) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	r := img.Bounds()
	w := r.Dx()
	h := r.Dy()

	if options.MinWidth > 0 && w < options.MinWidth {
		return dataTooSmall
	}
	if options.MinHeight > 0 && h < options.MinHeight {
		return dataTooSmall
	}
	if options.MaxWidth > 0 && w > options.MaxWidth {
		return dataTooLarge
	}
	if options.MaxHeight > 0 && h > options.MaxHeight {
		return dataTooLarge
	}

	return nil
}

func persistent(path string, ctx context, options Options) error {
	lumber.Start(ctx.dict())

	defer lumber.Finish(ctx.dict())

	return os.Rename(path, ctx.path)
}
