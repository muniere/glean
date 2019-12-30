package action

import (
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Status(w *rpc.Gateway, ctx *scope.Context) error {
	return w.Success(ctx.Queue.List())
}
