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
type RequestHook func(*rpc.Request)
type ResponseHook func(*rpc.Request)
type ErrorHook func(error)

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

	x.OnRequest(func(request *rpc.Request) {
		lumber.Info(box.Dict{"module": "producer", "event": "request", "command": request.Action})
	})

	x.OnError(rpc.Accept, func(err error) {
		if rpc.IsClosedConn(err) {
			lumber.Trace(box.Dict{"module": "producer", "event": "abort", "error": err.Error()})
			return
		}

		lumber.Error(box.Dict{"module": "producer", "event": "error", "error": err.Error()})
	})

	x.OnError(rpc.Handle, func(err error) {
		if rpc.IsTimeout(err) {
			lumber.Trace(box.Dict{"module": "producer", "event": "poll.timeout"})
			return
		}

		if rpc.IsEOF(err) {
			lumber.Trace(box.Dict{"module": "producer", "event": "poll.eof"})
			return
		}

		lumber.Error(box.Dict{"module": "producer", "event": "error", "error": err.Error()})
	})

	return x
}

func (x *Producer) Start() error {
	lumber.Info(box.Dict{"module": "producer", "event": "start"})
	return x.daemon.Start()
}

func (x *Producer) Stop() error {
	lumber.Info(box.Dict{"module": "producer", "event": "stop"})
	return x.daemon.Stop()
}

func (x *Producer) Wait() {
	x.daemon.Wait()
}

func (x *Producer) OnRequest(hook RequestHook) {
	x.daemon.OnRequest(func(request *rpc.Request) {
		hook(request)
	})
}

func (x *Producer) OnResponse(hook ResponseHook) {
	x.daemon.OnResponse(func(request *rpc.Request) {
		hook(request)
	})
}

func (x *Producer) OnError(phase rpc.Phase, hook ErrorHook) {
	x.daemon.OnError(phase, func(err error) {
		hook(err)
	})
}

func (x *Producer) Register(key string, proc Proc) {
	x.daemon.Register(key, func(request *rpc.Request, gateway *rpc.Gateway) error {
		return proc(action.NewContext(request, gateway, x.queue))
	})
}

func (x *Producer) RegisterDefault(proc Proc) {
	x.daemon.RegisterDefault(func(request *rpc.Request, gateway *rpc.Gateway) error {
		return proc(action.NewContext(request, gateway, x.queue))
	})
}
