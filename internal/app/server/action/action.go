package action

import (
	"github.com/muniere/glean/internal/app/server/action/cancel"
	"github.com/muniere/glean/internal/app/server/action/clutch"
	"github.com/muniere/glean/internal/app/server/action/fallback"
	"github.com/muniere/glean/internal/app/server/action/scrape"
	"github.com/muniere/glean/internal/app/server/action/shared"
	"github.com/muniere/glean/internal/app/server/action/status"
)

type (
	Context = shared.Context
)

var (
	NewContext = shared.NewContext
	Cancel     = cancel.Perform
	Clutch     = clutch.Perform
	Scrape     = scrape.Perform
	Status     = status.Perform
	Fallback   = fallback.Perform
)
