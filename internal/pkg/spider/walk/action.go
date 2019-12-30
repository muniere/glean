package walk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/spider/log"
)

type command struct {
	uri *url.URL
}

func Perform(uri *url.URL, options Options) ([]*url.URL, error) {
	cmd := command{uri: uri}
	ctx := compose(cmd, options)

	data, err := fetch(ctx, options)
	if err != nil {
		return nil, err
	}

	return scrape(data, ctx, options)
}

func compose(cmd command, options Options) context {
	return context{uri: cmd.uri}
}

func fetch(context context, options Options) (json.RawMessage, error) {
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

	return ioutil.ReadAll(res.Body)
}

func scrape(data json.RawMessage, context context, options Options) ([]*url.URL, error) {
	log.Debug("start", "scrape", context.dict())

	defer log.Debug("finish", "scrape", context.dict())

	re := regexp.MustCompile(".*\\.(jpg|png|gif)")

	var uris []*url.URL

	err := jsonic.Walk(data, func(val interface{}) error {
		switch v := val.(type) {
		case string:
			if !re.MatchString(v) {
				return nil
			}

			u, err := url.Parse(v)
			if err != nil {
				return err
			}
			if len(u.Scheme) == 0 || len(u.Host) == 0 {
				return nil
			}

			uris = append(uris, u)
			return nil

		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return uris, nil
}
