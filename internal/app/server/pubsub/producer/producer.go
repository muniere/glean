package producer

import (
	pubsub "github.com/muniere/glean/internal/app/server/pubsub/axiom"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/rpc"
	. "github.com/muniere/glean/internal/pkg/stdlib"
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

func NewProducer(queue *task.Queue, config Config) *Producer {
	x := &Producer{
		daemon: rpc.NewDaemon(config.Address, config.Port),
		queue:  queue,
	}

	x.RegisterRequestHandler(func(request *rpc.Request) error {
		lumber.Info(Dict{"module": "producer", "event": "request", "command": request.Action})
		return nil
	})

	x.RegisterErrorHandler(rpc.Accept, func(err error) {
		if rpc.IsClosedConn(err) {
			lumber.Trace(Dict{"module": "producer", "event": "abort", "error": err.Error()})
			return
		}

		lumber.Error(Dict{"module": "producer", "event": "error", "error": err.Error()})
	})

	x.RegisterErrorHandler(rpc.Handle, func(err error) {
		if rpc.IsTimeout(err) {
			lumber.Trace(Dict{"module": "producer", "event": "poll.timeout"})
			return
		}

		if rpc.IsEOF(err) {
			lumber.Trace(Dict{"module": "producer", "event": "poll.eof"})
			return
		}

		lumber.Error(Dict{"module": "producer", "event": "error", "error": err.Error()})
	})

	return x
}

func (x *Producer) Start() error {
	lumber.Info(Dict{"module": "producer", "event": "start"})
	return x.daemon.Start()
}

func (x *Producer) Stop() error {
	lumber.Info(Dict{"module": "producer", "event": "stop"})
	return x.daemon.Stop()
}

func (x *Producer) Wait() {
	x.daemon.Wait()
}

func (x *Producer) RegisterRequestHook(hook *pubsub.RequestHook) {
	x.daemon.OnRequest(&requestHookAdapter{hook})
}

func (x *Producer) RegisterRequestHandler(handler func(req *rpc.Request) error) {
	x.daemon.OnRequest(&requestHookAdapter{&pubsub.RequestHook{Handler: handler}})
}

func (x *Producer) RegisterResponseHook(hook *pubsub.ResponseHook) {
	x.daemon.OnResponse(&responseHookAdapter{hook})
}

func (x *Producer) RegisterResponseHandler(handler func(req *rpc.Request) error) {
	x.daemon.OnResponse(&responseHookAdapter{&pubsub.ResponseHook{Handler: handler}})
}

func (x *Producer) RegisterErrorHook(phase rpc.Phase, hook *pubsub.ErrorHook) {
	x.daemon.OnError(phase, &errorHookAdapter{hook})
}

func (x *Producer) RegisterErrorHandler(phase rpc.Phase, handler func(err error)) {
	x.daemon.OnError(phase, &errorHookAdapter{&pubsub.ErrorHook{Handler: handler}})
}

func (x *Producer) RegisterAction(key string, action *pubsub.Action) {
	x.daemon.Register(key, &actionAdapter{action, x.queue})
}

func (x *Producer) RegisterHandler(key string, handler func(ctx *pubsub.Context) error) {
	x.daemon.Register(key, &actionAdapter{&pubsub.Action{Handler: handler}, x.queue})
}

func (x *Producer) RegisterDefaultAction(action *pubsub.Action) {
	x.daemon.RegisterDefault(&actionAdapter{action, x.queue})
}

func (x *Producer) RegisterDefaultHandler(handler func(ctx *pubsub.Context) error) {
	x.daemon.RegisterDefault(&actionAdapter{&pubsub.Action{Handler: handler}, x.queue})
}
