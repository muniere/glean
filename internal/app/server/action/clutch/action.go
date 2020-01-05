package clutch

import (
	"github.com/muniere/glean/internal/app/server/action/shared"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Perform(ctx *shared.Context) error {
	var payload rpc.ClutchPayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Enqueue(rpc.Clutch, payload.URI, payload.Prefix)
	if err != nil {
		return ctx.Gateway.Error(box.Failure{
			Message: err.Error(),
		})
	}

	lumber.Info(box.Dict{
		"module": "producer",
		"event":  "job::produce",
		"job":    job,
	})

	return ctx.Gateway.Success(job)
}
