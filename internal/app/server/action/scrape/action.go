package scrape

import (
	pubsub "github.com/muniere/glean/internal/app/server/pubsub/axiom"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func NewAction() *pubsub.Action {
	return &pubsub.Action{
		Handler: perform,
	}
}

func perform(ctx *pubsub.Context) error {
	var payload rpc.ScrapePayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Enqueue(rpc.Scrape, payload.URI, payload.Prefix)
	if err != nil {
		return ctx.Gateway.Error(box.Failure{
			Message: err.Error(),
		})
	}

	lumber.Info(box.Dict{"module": "producer", "event": "job::produce", "job": job})

	return ctx.Gateway.Success(job)
}
