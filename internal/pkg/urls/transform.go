package urls

import (
	"net/url"
)

func Unique(urls []*url.URL) []*url.URL {
	var arr []*url.URL
	var dict = map[string]bool{}

	for _, u := range urls {
		_, ok := dict[u.String()]
		if ok {
			continue
		} else {
			arr = append(arr, u)
			dict[u.String()] = true
		}
	}

	return arr
}
