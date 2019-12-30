package action

import (
	"github.com/muniere/glean/internal/app/server/action/cancel"
	"github.com/muniere/glean/internal/app/server/action/clutch"
	"github.com/muniere/glean/internal/app/server/action/context"
	"github.com/muniere/glean/internal/app/server/action/fallback"
	"github.com/muniere/glean/internal/app/server/action/scrape"
	"github.com/muniere/glean/internal/app/server/action/status"
)

type (
	Context = context.Context
)

var (
	NewContext = context.NewContext
	Cancel     = cancel.Perform
	Clutch     = clutch.Perform
	Scrape     = scrape.Perform
	Status     = status.Perform
	Fallback   = fallback.Perform
)
