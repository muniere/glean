package context

import (
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Context struct {
	Request *rpc.Request
	Gateway *rpc.Gateway
	Queue   *task.Queue
}

func NewContext(request *rpc.Request, gateway *rpc.Gateway, queue *task.Queue) *Context {
	return &Context{
		Request: request,
		Gateway: gateway,
		Queue:   queue,
	}
}
