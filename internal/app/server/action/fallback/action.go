package fallback

import (
	"github.com/muniere/glean/internal/app/server/action/shared"
)

func Perform(ctx *shared.Context) error {
	return ctx.Gateway.Error(nil)
}
