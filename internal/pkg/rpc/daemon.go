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
)

var ConnectionClosed = errors.New("connection closed")

type Daemon struct {
	Address      string
	Port         int
	Window       int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	delegate     net.Listener
	procs        map[string]Proc
	fallback     Proc
	group        *sync.WaitGroup
}

type Proc func(net.Conn, []byte) error

func NewDaemon(addr string, port int) *Daemon {
	return &Daemon{
		Address:      addr,
		Port:         port,
		Window:       1024,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		delegate:     nil,
		procs:        map[string]Proc{},
		fallback:     nil,
		group:        &sync.WaitGroup{},
	}
}

func (d *Daemon) Register(key string, proc Proc) {
	d.procs[key] = proc
}

func (d *Daemon) Unregister(key string) {
	delete(d.procs, key)
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

	d.delegate = listener
	d.group.Add(1)

	go func() {
		defer d.group.Done()

		for {
			err := d.poll()

			if err == nil {
				continue
			}
			if err == ConnectionClosed {
				log.Trace(jsonic.MustEncode(box.Dict{
					"module":  "daemon",
					"action":  "accept",
					"label":   "stop",
					"message": err.Error(),
				}))
				break
			}

			log.Warn(err)
			break
		}
	}()

	return nil
}

func (d *Daemon) Stop() error {
	return d.delegate.Close()
}

func (d *Daemon) Wait() {
	d.group.Wait()
}

func (d *Daemon) poll() error {
	con, err := d.delegate.Accept()
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
			log.Trace(jsonic.MustEncode(box.Dict{
				"module": "daemon",
				"action": "poll",
				"label":  "Timeout",
			}))
			return
		}
		if err == io.EOF {
			log.Trace(jsonic.MustEncode(box.Dict{
				"module": "daemon",
				"action": "poll",
				"label":  "EOF",
			}))
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

	var msg Request
	if err := jsonic.Unmarshal(req, &msg); err != nil {
		return err
	}

	// response
	proc, ok := d.procs[msg.Action]
	if ok {
		return proc(con, req)
	}
	if d.fallback != nil {
		return d.fallback(con, req)
	}
	return errors.New(fmt.Sprintf("unsupported type: %v", msg.Action))
}
