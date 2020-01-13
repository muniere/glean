package status

import (
	pubsub "github.com/muniere/glean/internal/app/server/pubsub/base"
)

func NewAction() *pubsub.Action {
	return &pubsub.Action{
		Handler: perform,
	}
}

func perform(ctx *pubsub.Context) error {
	return ctx.Gateway.Success(ctx.Queue.List())
}
