package server

import (
	"encoding/json"
	"net"

	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type Server struct {
	delegate *rpc.Daemon
	jobs     *task.Queue
}

type Proc func(*relay.Gateway, *scope.Context) error

func New(addr string, port int) *Server {
	s := &Server{
		delegate: rpc.NewDaemon(addr, port),
		jobs:     task.NewQueue(),
	}

	s.Register("status", action.Status)
	s.Register("launch", action.Launch)
	s.Register("cancel", action.Cancel)
	s.RegisterDefault(action.Uncaught)

	return s
}

func (s *Server) Start() error {
	return s.delegate.Start()
}

func (s *Server) Stop() error {
	return s.delegate.Stop()
}

func (s *Server) Register(key string, proc Proc) {
	s.delegate.Register(key, func(con net.Conn, req []byte) error {
		var r rpc.Request
		if err := json.Unmarshal(req, &r); err != nil {
			return err
		}

		return proc(
			relay.NewGateway(con),
			scope.NewContext(&r, s.jobs),
		)
	})
}

func (s *Server) RegisterDefault(proc Proc) {
	s.delegate.RegisterDefault(func(con net.Conn, req []byte) error {
		var r rpc.Request
		if err := json.Unmarshal(req, &r); err != nil {
			return err
		}
		return proc(
			relay.NewGateway(con),
			scope.NewContext(&r, s.jobs),
		)
	})
}
