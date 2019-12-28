package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/muniere/glean/internal/pkg/chars"
	"github.com/muniere/glean/internal/pkg/packet"
)

const timeout = 30 * time.Second

type Agent struct {
	Address      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(address string, port int) *Agent {
	return &Agent{
		Address:      address,
		Port:         port,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
}

func (c *Agent) Dial() (net.Conn, error) {
	addr := fmt.Sprintf("%v:%d", c.Address, c.Port)
	return net.Dial("tcp", addr)
}

func (c *Agent) Submit(request packet.Request) ([]byte, error) {
	// connect
	con, err := c.Dial()
	if err != nil {
		return nil, err
	}

	// request
	buf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req := append(buf, chars.NUL)
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

	return res, nil
}
