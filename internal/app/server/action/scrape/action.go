package scrape

import (
	pubsub "github.com/muniere/glean/internal/app/server/pubsub/axiom"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
	. "github.com/muniere/glean/internal/pkg/stdlib"
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
		return ctx.Gateway.Error(pubsub.Failure{
			Message: err.Error(),
		})
	}

	lumber.Info(NewDict(Pair("module", "producer"), Pair("event", "job::produce"), Pair("job", job)))

	return ctx.Gateway.Success(job)
}
