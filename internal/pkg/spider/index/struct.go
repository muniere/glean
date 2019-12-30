package index

import (
	"net/url"
)

type SiteInfo struct {
	URI   *url.URL
	Title string
	Links []*url.URL
}
