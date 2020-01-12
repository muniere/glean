package clutch

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
	var payload rpc.ClutchPayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Enqueue(rpc.Clutch, payload.URI, payload.Prefix)
	if err != nil {
		return ctx.Gateway.Error(pubsub.Failure{
			Message: err.Error(),
		})
	}

	lumber.Info(Dict{"module": "producer", "event": "job::produce", "job": job})

	return ctx.Gateway.Success(job)
}
