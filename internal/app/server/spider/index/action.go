package index

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"gopkg.in/xmlpath.v2"

	"github.com/muniere/glean/internal/app/server/spider/log"
)

type command struct {
	uri *url.URL
}

func Perform(uri *url.URL, options Options) (*SiteInfo, error) {
	cmd := command{uri: uri}
	ctx := compose(cmd, options)

	doc, err := fetch(ctx, options)
	if err != nil {
		return nil, err
	}

	return scrape(doc, ctx, options)
}

func compose(cmd command, options Options) context {
	return context{uri: cmd.uri}
}

func fetch(context context, options Options) (*xmlpath.Node, error) {
	log.Debug("start", "fetch", context.dict())

	defer log.Debug("finish", "fetch", context.dict())

	res, err := http.Get(context.uri.String())
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	r := bufio.NewReader(res.Body)
	data, err := r.Peek(1024)
	if err != nil {
		return nil, err
	}

	enc, _, ok := charset.DetermineEncoding(data, res.Header.Get("Content-Type"))
	if ok {
		return xmlpath.ParseHTML(enc.NewDecoder().Reader(r))
	} else {
		return xmlpath.ParseHTML(res.Body)
	}
}

func scrape(doc *xmlpath.Node, context context, options Options) (*SiteInfo, error) {
	log.Debug("start", "scrape", context.dict())

	defer log.Debug("finish", "scrape", context.dict())

	title := scrapeTitle(doc)

	re := regexp.MustCompile(".*\\.(jpg|png|gif)")

	hrefs, err := scrapeURLs(doc, "//a/@href", re, options.Grep)
	if err != nil {
		return nil, err
	}

	srcs, err := scrapeURLs(doc, "//img/@src", re, options.Grep)
	if err != nil {
		return nil, err
	}

	links := append(hrefs, srcs...)

	log.Result(len(links), context.dict())

	info := SiteInfo{
		URI:   context.uri,
		Title: title,
		Links: links,
	}

	return &info, nil
}

func scrapeTitle(doc *xmlpath.Node) string {
	xpath := xmlpath.MustCompile("//title")
	iter := xpath.Iter(doc)

	if iter.Next() {
		return iter.Node().String()
	} else {
		return ""
	}
}

func scrapeURLs(doc *xmlpath.Node, path string, pattern *regexp.Regexp, grep *regexp.Regexp) ([]*url.URL, error) {
	var result []*url.URL

	xpath := xmlpath.MustCompile(path)
	iter := xpath.Iter(doc)

	for iter.Next() {
		val := iter.Node().String()

		if pattern != nil && !pattern.MatchString(val) {
			continue
		}
		if grep != nil && !grep.MatchString(val) {
			continue
		}

		s := strings.Replace(val, " ", "+", -1)
		u, err := url.Parse(s)
		if err != nil {
			continue
		}

		result = append(result, u)
	}

	return result, nil
}
