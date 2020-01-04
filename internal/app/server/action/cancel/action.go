package cancel

import (
	"github.com/muniere/glean/internal/app/server/action/shared"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Perform(ctx *shared.Context) error {
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

	lumber.Info(box.Dict{
		"module": "producer",
		"action": "cancel",
		"job":    job,
	})

	return ctx.Gateway.Success(job)
}
