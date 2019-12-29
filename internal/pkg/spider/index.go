package spider

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
	"gopkg.in/xmlpath.v2"

	"github.com/muniere/glean/internal/pkg/jsonic"
)

type IndexOptions struct {
	Grep *regexp.Regexp
}

func Index(uri *url.URL, options IndexOptions) (*SiteInfo, error) {
	log.Debug(jsonic.MustEncode(dict{
		"label":     "start",
		"action":    "fetch",
		"param.uri": uri.String(),
	}))

	doc, err := fetch(uri)
	if err != nil {
		return nil, err
	}

	log.Debug(jsonic.MustEncode(dict{
		"label":     "start",
		"action":    "scrape",
		"param.uri": uri.String(),
	}))

	title := title(doc)

	re := regexp.MustCompile(".*\\.(jpg|png|gif)")

	hrefs, err := scrape(doc, "//a/@href", re, options)
	if err != nil {
		return nil, err
	}

	srcs, err := scrape(doc, "//img/@src", re, options)
	if err != nil {
		return nil, err
	}

	links := append(hrefs, srcs...)

	log.Debugf(jsonic.MustEncode(dict{
		"label":     "result",
		"result":    len(links),
		"param.uri": uri.String(),
	}))

	info := SiteInfo{
		URI:   uri,
		Title: title,
		Links: links,
	}

	return &info, nil
}

func fetch(uri *url.URL) (*xmlpath.Node, error) {
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

func title(doc *xmlpath.Node) string {
	xpath := xmlpath.MustCompile("//title")
	iter := xpath.Iter(doc)

	if iter.Next() {
		return iter.Node().String()
	} else {
		return ""
	}
}

func scrape(doc *xmlpath.Node, path string, pattern *regexp.Regexp, options IndexOptions) ([]*url.URL, error) {
	var result []*url.URL

	xpath := xmlpath.MustCompile(path)
	iter := xpath.Iter(doc)

	for iter.Next() {
		val := iter.Node().String()

		if pattern != nil && !pattern.MatchString(val) {
			continue
		}
		if options.Grep != nil && !options.Grep.MatchString(val) {
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
