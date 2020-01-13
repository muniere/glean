package base

import (
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Action struct {
	Handler func(ctx *Context) error
}

type Context struct {
	Request *rpc.Request
	Gateway *rpc.Gateway
	Queue   *task.Queue
}

type RequestHook struct {
	Handler func(*rpc.Request) error
}

type ResponseHook struct {
	Handler func(*rpc.Request) error
}

type ErrorHook struct {
	Handler func(error)
}

