package action

import (
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Launch(w *relay.Gateway, ctx *scope.Context) error {
	var payload rpc.LaunchPayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Enqueue(payload.URI)
	if err != nil {
		return w.Error(box.Failure{
			Message: err.Error(),
		})
	}

	return w.Success(job)
}
