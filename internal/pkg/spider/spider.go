package spider

import (
	"net/url"

	"github.com/muniere/glean/internal/pkg/spider/download"
	"github.com/muniere/glean/internal/pkg/spider/index"
	"github.com/muniere/glean/internal/pkg/spider/walk"
)

type SiteInfo = index.SiteInfo
type WalkOptions = walk.Options
type IndexOptions = index.Options
type DownloadOptions = download.Options

func Walk(url *url.URL, options WalkOptions) ([]*url.URL, error) {
	return walk.Perform(url, options)
}

func Index(url *url.URL, options IndexOptions) (*SiteInfo, error) {
	return index.Perform(url, options)
}

func Download(urls []*url.URL, options DownloadOptions) error {
	return download.Perform(urls, options)
}
