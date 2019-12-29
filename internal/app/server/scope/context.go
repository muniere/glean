package scope

import (
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Context struct {
	Request *rpc.Request
	Queue   *task.Queue
}

func NewContext(req *rpc.Request, queue *task.Queue) *Context {
	return &Context{
		Request: req,
		Queue:   queue,
	}
}
