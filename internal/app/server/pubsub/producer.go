package pubsub

import (
	"net"

	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Producer struct {
	daemon *rpc.Daemon
	queue  *task.Queue
}

type ProducerConfig struct {
	Address string
	Port    int
}

type Proc func(*action.Context) error

func NewProducer(queue *task.Queue, config ProducerConfig) *Producer {
	s := &Producer{
		daemon: rpc.NewDaemon(config.Address, config.Port),
		queue:  queue,
	}

	s.Register(rpc.Status, action.Status)
	s.Register(rpc.Scrape, action.Scrape)
	s.Register(rpc.Clutch, action.Clutch)
	s.Register(rpc.Cancel, action.Cancel)
	s.RegisterDefault(action.Fallback)

	return s
}

func (p *Producer) Start() error {
	return p.daemon.Start()
}

func (p *Producer) Stop() error {
	return p.daemon.Stop()
}

func (p *Producer) Wait() {
	p.daemon.Wait()
}

func (p *Producer) Register(key string, proc Proc) {
	p.daemon.Register(key, func(con net.Conn, req []byte) error {
		var request rpc.Request
		if err := jsonic.Unmarshal(req, &request); err != nil {
			return err
		}

		gateway := rpc.NewGateway(con)
		context := action.NewContext(&request, gateway, p.queue)

		return proc(context)
	})
}

func (p *Producer) RegisterDefault(proc Proc) {
	p.daemon.RegisterDefault(func(con net.Conn, req []byte) error {
		var request rpc.Request
		if err := jsonic.Unmarshal(req, &request); err != nil {
			return err
		}

		gateway := rpc.NewGateway(con)
		context := action.NewContext(&request, gateway, p.queue)

		return proc(context)
	})
}
