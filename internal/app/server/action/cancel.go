package action

import (
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Cancel(w *relay.Gateway, ctx *scope.Context) error {
	var payload rpc.CancelPayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Remove(payload.ID)
	if err != nil {
		return w.Error(box.Failure{
			Message: err.Error(),
		})
	}

	return w.Success(job)
}
