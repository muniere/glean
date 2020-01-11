package producer

import (
	. "github.com/muniere/glean/internal/app/server/pubsub/axiom"
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

type actionAdapter struct {
	action *Action
	queue  *task.Queue
}

func (w *actionAdapter) Invoke(req *rpc.Request, gateway *rpc.Gateway) error {
	return w.action.Handler(
		&Context{
			Request: req,
			Gateway: gateway,
			Queue: w.queue,
		},
	)
}

type requestHookAdapter struct {
	hook *RequestHook
}

func (w *requestHookAdapter) Invoke(req *rpc.Request) error {
	return w.hook.Handler(req)
}

type responseHookAdapter struct {
	hook *ResponseHook
}

func (w *responseHookAdapter) Invoke(req *rpc.Request) error {
	return w.hook.Handler(req)
}

type errorHookAdapter struct {
	hook *ErrorHook
}

func (w *errorHookAdapter) Invoke(err error) {
	w.hook.Handler(err)
}

func NewProducer(queue *task.Queue, config Config) *Producer {
	x := &Producer{
		daemon: rpc.NewDaemon(config.Address, config.Port),
		queue:  queue,
	}

	x.RegisterRequestHandler(func(request *rpc.Request) error {
		lumber.Info(box.Dict{"module": "producer", "event": "request", "command": request.Action})
		return nil
	})

	x.RegisterErrorHandler(rpc.Accept, func(err error) {
		if rpc.IsClosedConn(err) {
			lumber.Trace(box.Dict{"module": "producer", "event": "abort", "error": err.Error()})
			return
		}

		lumber.Error(box.Dict{"module": "producer", "event": "error", "error": err.Error()})
	})

	x.RegisterErrorHandler(rpc.Handle, func(err error) {
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

func (x *Producer) RegisterRequestHook(hook *RequestHook) {
	x.daemon.OnRequest(&requestHookAdapter{hook})
}

func (x *Producer) RegisterRequestHandler(handler func(req *rpc.Request) error) {
	x.daemon.OnRequest(&requestHookAdapter{&RequestHook{Handler: handler}})
}

func (x *Producer) RegisterResponseHook(hook *ResponseHook) {
	x.daemon.OnResponse(&responseHookAdapter{hook})
}

func (x *Producer) RegisterResponseHandler(handler func(req *rpc.Request) error) {
	x.daemon.OnResponse(&responseHookAdapter{&ResponseHook{Handler: handler}})
}

func (x *Producer) RegisterErrorHook(phase rpc.Phase, hook *ErrorHook) {
	x.daemon.OnError(phase, &errorHookAdapter{hook})
}

func (x *Producer) RegisterErrorHandler(phase rpc.Phase, handler func(err error)) {
	x.daemon.OnError(phase, &errorHookAdapter{&ErrorHook{Handler: handler}})
}

func (x *Producer) RegisterAction(key string, action *Action) {
	x.daemon.Register(key, &actionAdapter{action, x.queue})
}

func (x *Producer) RegisterHandler(key string, handler func(ctx *Context) error) {
	x.daemon.Register(key, &actionAdapter{&Action{Handler: handler}, x.queue})
}

func (x *Producer) RegisterDefaultAction(action *Action) {
	x.daemon.RegisterDefault(&actionAdapter{action, x.queue})
}

func (x *Producer) RegisterDefaultHandler(handler func(ctx *Context) error) {
	x.daemon.RegisterDefault(&actionAdapter{&Action{Handler: handler}, x.queue})
}
