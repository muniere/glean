package action

import (
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/app/server/scope"
)

func Status(w *relay.Gateway, ctx *scope.Context) error {
	return w.Success(ctx.Jobs.List())
}
