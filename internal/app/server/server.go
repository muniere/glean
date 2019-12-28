package server

import (
	"encoding/json"
	"net"

	"github.com/muniere/glean/internal/app/server/action"
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/pkg/daemon"
	"github.com/muniere/glean/internal/pkg/packet"
)

type Server struct {
	delegate *daemon.Daemon
}

func New(addr string, port int) *Server {
	s := &Server{
		delegate: daemon.New(addr, port),
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

func (s *Server) Register(key string, proc func(*relay.Gateway, *packet.Request) error) {
	s.delegate.Register(key, func(con net.Conn, req []byte) error {
		var request packet.Request
		if err := json.Unmarshal(req, &request); err != nil {
			return err
		}
		return proc(relay.NewGateway(con), &request)
	})
}

func (s *Server) RegisterDefault(proc func(*relay.Gateway, *packet.Request) error) {
	s.delegate.RegisterDefault(func(con net.Conn, req []byte) error {
		var request packet.Request
		if err := json.Unmarshal(req, &request); err != nil {
			return err
		}
		return proc(relay.NewGateway(con), &request)
	})
}
