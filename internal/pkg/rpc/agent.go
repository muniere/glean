package rpc

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/muniere/glean/internal/pkg/ascii"
	"github.com/muniere/glean/internal/pkg/jsonic"
)

type Agent struct {
	Address      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewAgent(address string, port int) *Agent {
	return &Agent{
		Address:      address,
		Port:         port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (c *Agent) Dial() (net.Conn, error) {
	addr := fmt.Sprintf("%v:%d", c.Address, c.Port)
	return net.Dial("tcp", addr)
}

func (c *Agent) Submit(request *Request) (*Response, error) {
	// connect
	con, err := c.Dial()
	if err != nil {
		return nil, err
	}

	// request
	buf, err := jsonic.Marshal(request)
	if err != nil {
		return nil, err
	}
	req := append(buf, ascii.NUL)
	_ = con.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	_, err = con.Write(req)
	if err != nil {
		return nil, err
	}

	// response
	_ = con.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	res, err := ioutil.ReadAll(con)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := jsonic.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
