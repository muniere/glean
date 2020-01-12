package producer

import (
	pubsub "github.com/muniere/glean/internal/app/server/pubsub/axiom"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

//
// Adapter for rpc.Action
//
type actionAdapter struct {
	action *pubsub.Action
	queue  *task.Queue
}

func (w *actionAdapter) Invoke(req *rpc.Request, gateway *rpc.Gateway) error {
	return w.action.Handler(
		&pubsub.Context{
			Request: req,
			Gateway: gateway,
			Queue:   w.queue,
		},
	)
}

//
// Adapter for rpc.RequestHook
//
type requestHookAdapter struct {
	hook *pubsub.RequestHook
}

func (w *requestHookAdapter) Invoke(req *rpc.Request) error {
	return w.hook.Handler(req)
}

//
// Adapter for rpc.ResponseHook
//
type responseHookAdapter struct {
	hook *pubsub.ResponseHook
}

func (w *responseHookAdapter) Invoke(req *rpc.Request) error {
	return w.hook.Handler(req)
}

//
// Adapter for rpc.ErrorHook
//
type errorHookAdapter struct {
	hook *pubsub.ErrorHook
}

func (w *errorHookAdapter) Invoke(err error) {
	w.hook.Handler(err)
}
