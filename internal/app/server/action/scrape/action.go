package scrape

import (
	"github.com/muniere/glean/internal/app/server/action/shared"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Perform(ctx *shared.Context) error {
	var payload rpc.ScrapePayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Enqueue(rpc.Scrape, payload.URI)
	if err != nil {
		return ctx.Gateway.Error(box.Failure{
			Message: err.Error(),
		})
	}

	return ctx.Gateway.Success(job)
}