package fallback

import (
	"github.com/muniere/glean/internal/app/server/action/context"
)

func Perform(ctx *context.Context) error {
	return ctx.Gateway.Error(nil)
}
