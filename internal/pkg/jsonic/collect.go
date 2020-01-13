package jsonic

import (
	"encoding/json"
	"net/url"

	"github.com/muniere/glean/internal/pkg/urls"
)

func Collect(data json.RawMessage, test func(string) bool) ([]*url.URL, error) {
	var uris []*url.URL

	err := Walk(data, func(val interface{}) error {
		switch v := val.(type) {
		case string:
			if !test(v) {
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

	return urls.Unique(uris), nil
}
