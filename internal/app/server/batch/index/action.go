package index

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"

	"github.com/muniere/glean/internal/app/server/batch/lumber"
	"github.com/muniere/glean/internal/pkg/std"
	"github.com/muniere/glean/internal/pkg/xml"
)

//
// Struct
//
type Options struct {
	Grep *regexp.Regexp
}

type SiteInfo struct {
	URI   *url.URL
	Title string
	Links []*url.URL
}

type docInfo struct {
	title string
	links []*url.URL
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
func Perform(uri *url.URL, options Options) (*SiteInfo, error) {
	cmd := command{uri: uri}
	ctx := compose(cmd, options)

	doc, err := fetch(ctx, options)
	if err != nil {
		return nil, err
	}

	info, err := scrape(doc, ctx, options)
	if err != nil {
		return nil, err
	}

	return bundle(info, ctx, options)
}

func compose(cmd command, options Options) context {
	return context{uri: cmd.uri}
}

func fetch(ctx context, options Options) (*xml.Node, error) {
	lumber.Start(ctx.dict())
	defer lumber.Finish(ctx.dict())
	return doFetch(ctx.uri)
}

func scrape(doc *xml.Node, ctx context, options Options) (*docInfo, error) {
	lumber.Start(ctx.dict())
	defer lumber.Finish(ctx.dict())
	return doScrape(doc, options.Grep)
}

func bundle(info *docInfo, ctx context, options Options) (*SiteInfo, error) {
	lumber.Result(len(info.links), ctx.dict())
	return &SiteInfo{
		ctx.uri,
		info.title,
		info.links,
	}, nil
}

//
// Helper
//
func doFetch(uri *url.URL) (*xml.Node, error) {
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

	return xml.Parse(res)
}

func doScrape(doc *xml.Node, grep *regexp.Regexp) (*docInfo, error) {
	pattern := regexp.MustCompile(".*\\.(jpg|png|gif)")

	title := xml.Title(doc)
	links := xml.Collect(doc, func(node *xml.Node) bool {
		if pattern != nil && !pattern.Match(node.Bytes()) {
			return false
		}
		if grep != nil && !grep.Match(node.Bytes()) {
			return false
		}
		return true
	})

	return &docInfo{
		title: title,
		links: links,
	}, nil
}
