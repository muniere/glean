package relay

import (
	"net"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

type Gateway struct {
	delegate net.Conn
}

func NewGateway(con net.Conn) *Gateway {
	return &Gateway{
		delegate: con,
	}
}

func (w *Gateway) Success(payload interface{}) error {
	return w.Respond(rpc.Response{Ok: true, Payload: payload})
}

func (w *Gateway) Error(payload interface{}) error {
	return w.Respond(rpc.Response{Ok: false, Payload: payload})
}

func (w *Gateway) Respond(response rpc.Response) error {
	res, err := jsonic.Marshal(response)
	if err != nil {
		return err
	}

	_, err = w.delegate.Write(res)
	if err != nil {
		return err
	}

	return nil
}
