package action

import (
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Clutch(w *rpc.Gateway, ctx *scope.Context) error {
	var payload rpc.ClutchPayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Enqueue(rpc.Clutch, payload.URI)
	if err != nil {
		return w.Error(box.Failure{
			Message: err.Error(),
		})
	}

	return w.Success(job)
}
