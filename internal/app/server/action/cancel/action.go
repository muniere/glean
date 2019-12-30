package cancel

import (
	"github.com/muniere/glean/internal/app/server/action/context"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Perform(ctx *context.Context) error {
	var payload rpc.CancelPayload
	if err := jsonic.Transcode(ctx.Request.Payload, &payload); err != nil {
		return err
	}

	job, err := ctx.Queue.Remove(payload.ID)
	if err != nil {
		return ctx.Gateway.Error(box.Failure{
			Message: err.Error(),
		})
	}

	return ctx.Gateway.Success(job)
}
