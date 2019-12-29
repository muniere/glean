package rpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/muniere/glean/internal/pkg/ascii"
)

type Daemon struct {
	Address      string
	Port         int
	Window       int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	delegate     net.Listener
	procs        map[string]Proc
	fallback     Proc
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

	defer func() {
		_ = d.Stop()
	}()

	for {
		d.poll()
	}
}

func (d *Daemon) Stop() error {
	return d.delegate.Close()
}

func (d *Daemon) poll() {
	con, err := d.delegate.Accept()
	if err != nil {
		log.Fatal(err)
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
			fmt.Println("Timeout")
			return
		}
		if err == io.EOF {
			fmt.Println("End of File")
			return
		}
		log.Fatal(err)
	}()
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
	if err := json.Unmarshal(req, &msg); err != nil {
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
