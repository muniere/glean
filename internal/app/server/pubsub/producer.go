package pubsub

import (
	"net"

	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/app/server/scope"
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

type Proc func(*rpc.Gateway, *scope.Context) error

func NewProducer(queue *task.Queue, config ProducerConfig) *Producer {
	s := &Producer{
		daemon: rpc.NewDaemon(config.Address, config.Port),
		queue:  queue,
	}

	s.Register(rpc.Status, action.Status)
	s.Register(rpc.Scrape, action.Scrape)
	s.Register(rpc.Clutch, action.Clutch)
	s.Register(rpc.Cancel, action.Cancel)
	s.RegisterDefault(action.Uncaught)

	return s
}

func (s *Producer) Start() error {
	return s.daemon.Start()
}

func (s *Producer) Stop() error {
	return s.daemon.Stop()
}

func (s *Producer) Wait() {
	s.daemon.Wait()
}

func (s *Producer) Register(key string, proc Proc) {
	s.daemon.Register(key, func(con net.Conn, req []byte) error {
		var r rpc.Request
		if err := jsonic.Unmarshal(req, &r); err != nil {
			return err
		}

		return proc(
			rpc.NewGateway(con),
			scope.NewContext(&r, s.queue),
		)
	})
}

func (s *Producer) RegisterDefault(proc Proc) {
	s.daemon.RegisterDefault(func(con net.Conn, req []byte) error {
		var r rpc.Request
		if err := jsonic.Unmarshal(req, &r); err != nil {
			return err
		}
		return proc(
			rpc.NewGateway(con),
			scope.NewContext(&r, s.queue),
		)
	})
}
