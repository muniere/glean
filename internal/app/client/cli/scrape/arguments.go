package scrape

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func normalize(args []string) ([]*url.URL, error) {
	var uris []*url.URL
	var errs []string

	for _, arg := range args {
		uri, err := url.Parse(arg)
		if err == nil && len(uri.Scheme) > 0 && len(uri.Host) > 0 {
			uris = append(uris, uri)
		} else {
			errs = append(errs, arg)
		}
	}

	if len(errs) > 0 {
		arg := strings.Join(errs, ", ")
		msg := fmt.Sprintf("values must be valid URLs: %s", arg)
		return nil, errors.New(msg)
	}

	return uris, nil
}
