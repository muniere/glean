package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/ascii"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/lumber"
)

var ConnectionClosed = errors.New("connection closed")

type Daemon struct {
	Address      string
	Port         int
	Window       int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	listener     net.Listener
	actions      map[string]Proc
	fallback     Proc
	preHooks     []Hook
	postHooks    []Hook
	group        *sync.WaitGroup
}

type Hook func(*Request)
type Proc func(*Request, *Gateway) error

func NewDaemon(addr string, port int) *Daemon {
	return &Daemon{
		Address:      addr,
		Port:         port,
		Window:       1024,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		listener:     nil,
		actions:      map[string]Proc{},
		fallback:     nil,
		preHooks:     []Hook{},
		postHooks:    []Hook{},
		group:        &sync.WaitGroup{},
	}
}

func (d *Daemon) PreHook(hook Hook) {
	d.preHooks = append(d.preHooks, hook)
}

func (d *Daemon) PostHook(hook Hook) {
	d.postHooks = append(d.postHooks, hook)
}

func (d *Daemon) Register(key string, proc Proc) {
	d.actions[key] = proc
}

func (d *Daemon) Unregister(key string) {
	delete(d.actions, key)
}

func (d *Daemon) RegisterDefault(proc Proc) {
	d.fallback = proc
}

func (d *Daemon) UnregisterDefault(proc Proc) {
	d.fallback = nil
}

func (d *Daemon) Start() error {
	address := fmt.Sprintf("%s:%d", d.Address, d.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	d.listener = listener
	d.group.Add(1)

	go func() {
		defer d.group.Done()

		for {
			err := d.poll()

			if err == nil {
				continue
			}
			if err == ConnectionClosed {
				lumber.Trace(box.Dict{
					"module":  "daemon",
					"action":  "abort",
					"message": err.Error(),
				})
				break
			}

			log.Warn(err)
			break
		}
	}()

	return nil
}

func (d *Daemon) Stop() error {
	return d.listener.Close()
}

func (d *Daemon) Wait() {
	d.group.Wait()
}

func (d *Daemon) poll() error {
	con, err := d.listener.Accept()
	if err != nil {
		// see https://github.com/golang/net/blob/master/http2/server.go#L676-L679
		if strings.Contains(err.Error(), "use of closed network connection") {
			return ConnectionClosed
		} else {
			return err
		}
	}

	go func() {
		defer func() {
			_ = con.Close()
		}()
		err := d.handle(con)
		if err == nil {
			return
		}
		e, ok := err.(net.Error)
		if ok && e.Timeout() {
			lumber.Trace(box.Dict{
				"module": "daemon",
				"action": "poll.timeout",
			})
			return
		}
		if err == io.EOF {
			lumber.Trace(box.Dict{
				"module": "daemon",
				"action": "poll.eof",
			})
			return
		}
		log.Error(err)
	}()

	return nil
}

func (d *Daemon) handle(con net.Conn) error {
	// request
	r := bufio.NewReader(con)

	_ = con.SetReadDeadline(time.Now().Add(d.ReadTimeout))
	req, err := r.ReadBytes(ascii.NUL)
	if err != nil {
		return err
	}

	req = req[0 : len(req)-1]

	var request Request
	if err := jsonic.Unmarshal(req, &request); err != nil {
		return err
	}

	// response
	for _, hook := range d.preHooks {
		hook(&request)
	}

	if err := d.perform(&request, NewGateway(con)); err != nil {
		return err
	}

	for _, hook := range d.postHooks {
		hook(&request)
	}

	return nil
}

func (d *Daemon) perform(request *Request, gateway *Gateway) error {
	action, ok := d.actions[request.Action]

	if ok {
		return action(request, gateway)
	}

	if d.fallback != nil {
		return d.fallback(request, gateway)
	}

	return errors.New(fmt.Sprintf("unsupported aciton: %s", request.Action))
}
