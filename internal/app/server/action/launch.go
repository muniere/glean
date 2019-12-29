package action

import (
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Launch(w *relay.Gateway, ctx *scope.Context) error {
	var payload rpc.LaunchPayload
	if err := ctx.Request.DecodePayload(&payload); err != nil {
		return err
	}

	if err := ctx.Queue.Enqueue(payload.Query); err != nil {
		f := box.Failure{
			Message: err.Error(),
		}
		return w.Error(f)
	}

	return w.Success(payload)
}
