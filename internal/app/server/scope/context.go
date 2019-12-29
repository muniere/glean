package scope

import (
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Context struct {
	Request *rpc.Request
	Jobs    *task.Queue
}

func NewContext(req *rpc.Request, jobs *task.Queue) *Context {
	return &Context{
		Request: req,
		Jobs:    jobs,
	}
}
