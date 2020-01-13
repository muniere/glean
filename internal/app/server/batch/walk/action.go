package walk

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/muniere/glean/internal/app/server/batch/lumber"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/std"
)

//
// Struct
//
type Options struct {
	Grep *regexp.Regexp
}

type command struct {
	uri *url.URL
}

type context struct {
	uri *url.URL
}

func (c *context) dict() std.Dict {
	return std.NewDict(
		std.Pair("uri", c.uri.String()),
	)
}

//
// Action
//
func Perform(uri *url.URL, options Options) ([]*url.URL, error) {
	cmd := command{uri: uri}
	ctx := compose(cmd, options)

	data, err := fetch(ctx, options)
	if err != nil {
		return nil, err
	}

	links, err := scrape(data, ctx, options)
	if err != nil {
		return nil, err
	}

	return bundle(links, ctx, options)
}

func compose(cmd command, options Options) context {
	return context{uri: cmd.uri}
}

func fetch(ctx context, options Options) ([]byte, error) {
	lumber.Start(ctx.dict())
	defer lumber.Finish(ctx.dict())
	return doFetch(ctx.uri)
}

func scrape(data []byte, ctx context, options Options) ([]*url.URL, error) {
	lumber.Start(ctx.dict())
	defer lumber.Finish(ctx.dict())
	return doScrape(data, options.Grep)
}

func bundle(links []*url.URL, ctx context, options Options) ([]*url.URL, error) {
	return links, nil
}

//
// Helper
//
func doFetch(uri *url.URL) ([]byte, error) {
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

	return ioutil.ReadAll(res.Body)
}

func doScrape(data []byte, grep *regexp.Regexp) ([]*url.URL, error) {
	pattern := regexp.MustCompile(".*\\.(jpg|png|gif)")

	links, err := jsonic.Collect(data, func(val string) bool {
		if pattern != nil && !pattern.MatchString(val) {
			return false
		}
		if grep != nil && !grep.MatchString(val) {
			return false
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return links, nil
}
