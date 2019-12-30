package spider

import (
	"github.com/muniere/glean/internal/app/server/spider/download"
	"github.com/muniere/glean/internal/app/server/spider/index"
	"github.com/muniere/glean/internal/app/server/spider/walk"
)

type (
	SiteInfo        = index.SiteInfo
	WalkOptions     = walk.Options
	IndexOptions    = index.Options
	DownloadOptions = download.Options
)

var (
	Walk     = walk.Perform
	Index    = index.Perform
	Download = download.Perform
)
