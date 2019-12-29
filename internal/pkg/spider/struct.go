package spider

import (
	"net/url"
)

type dict map[string]interface{}

type SiteInfo struct {
	URI   *url.URL
	Title string
	Links []*url.URL
}
