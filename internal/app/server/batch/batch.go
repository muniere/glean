package batch

import (
	"github.com/muniere/glean/internal/app/server/batch/download"
	"github.com/muniere/glean/internal/app/server/batch/index"
	"github.com/muniere/glean/internal/app/server/batch/walk"
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
