package fallback

import (
	pubsub "github.com/muniere/glean/internal/app/server/pubsub/axiom"
)

func NewAction() *pubsub.Action {
	return &pubsub.Action{
		Handler: func(ctx *pubsub.Context) error {
			return ctx.Gateway.Error(nil)
		},
	}
}
