package spider

import (
	"net/url"

	"github.com/muniere/glean/internal/pkg/spider/download"
	"github.com/muniere/glean/internal/pkg/spider/index"
)

type SiteInfo = index.SiteInfo
type IndexOptions = index.Options
type DownloadOptions = download.Options

func Index(url *url.URL, options IndexOptions) (*SiteInfo, error) {
	return index.Perform(url, options)
}

func Download(urls []*url.URL, options DownloadOptions) error {
	return download.Perform(urls, options)
}
