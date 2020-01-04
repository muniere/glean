package producer

import (
	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Producer struct {
	daemon *rpc.Daemon
	queue  *task.Queue
}

type Config struct {
	Address string
	Port    int
}

type Proc func(*action.Context) error

func NewProducer(queue *task.Queue, config Config) *Producer {
	x := &Producer{
		daemon: rpc.NewDaemon(config.Address, config.Port),
		queue:  queue,
	}

	x.Register(rpc.Status, action.Status)
	x.Register(rpc.Scrape, action.Scrape)
	x.Register(rpc.Clutch, action.Clutch)
	x.Register(rpc.Cancel, action.Cancel)
	x.RegisterDefault(action.Fallback)

	return x
}

func (x *Producer) Start() error {
	lumber.Info(box.Dict{
		"module": "producer",
		"action": "start",
	})
	return x.daemon.Start()
}

func (x *Producer) Stop() error {
	lumber.Info(box.Dict{
		"module": "producer",
		"action": "stop",
	})
	return x.daemon.Stop()
}

func (x *Producer) Wait() {
	x.daemon.Wait()
}

func (x *Producer) Register(key string, proc Proc) {
	x.daemon.Register(key, func(request *rpc.Request, gateway *rpc.Gateway) error {
		return proc(
			action.NewContext(request, gateway, x.queue),
		)
	})
}

func (x *Producer) RegisterDefault(proc Proc) {
	x.daemon.RegisterDefault(func(request *rpc.Request, gateway *rpc.Gateway) error {
		return proc(
			action.NewContext(request, gateway, x.queue),
		)
	})
}
