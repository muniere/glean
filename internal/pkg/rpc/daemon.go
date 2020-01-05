package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/muniere/glean/internal/pkg/ascii"
	"github.com/muniere/glean/internal/pkg/jsonic"
)

type Phase int

const (
	Accept = iota
	Handle
)

type Proc func(*Request, *Gateway) error
type RequestHook func(*Request)
type ResponseHook func(*Request)
type ErrorHook func(error)

type Daemon struct {
	Address       string
	Port          int
	Window        int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	listener      net.Listener
	actions       map[string]Proc
	fallback      Proc
	requestHooks  []RequestHook
	responseHooks []ResponseHook
	errorHooks    map[Phase][]ErrorHook
	group         *sync.WaitGroup
}

func NewDaemon(addr string, port int) *Daemon {
	return &Daemon{
		Address:       addr,
		Port:          port,
		Window:        1024,
		ReadTimeout:   1 * time.Second,
		WriteTimeout:  1 * time.Second,
		listener:      nil,
		actions:       map[string]Proc{},
		fallback:      nil,
		requestHooks:  []RequestHook{},
		responseHooks: []ResponseHook{},
		errorHooks:    map[Phase][]ErrorHook{},
		group:         &sync.WaitGroup{},
	}
}

func (d *Daemon) OnRequest(hook RequestHook) {
	d.requestHooks = append(d.requestHooks, hook)
}

func (d *Daemon) OnResponse(hook ResponseHook) {
	d.responseHooks = append(d.responseHooks, hook)
}

func (d *Daemon) OnError(phase Phase, hook ErrorHook) {
	d.errorHooks[phase] = append(d.subErrorHooks(phase), hook)
}

func (d *Daemon) subErrorHooks(phase Phase) []ErrorHook {
	arr, ok := d.errorHooks[phase]
	if ok {
		return arr
	} else {
		return []ErrorHook{}
	}
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

			for _, hook := range d.subErrorHooks(Accept) {
				hook(err)
			}
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
		return err
	}

	go func() {
		defer func() {
			_ = con.Close()
		}()

		err := d.handle(con)

		if err == nil {
			return
		}

		for _, hook := range d.subErrorHooks(Handle) {
			hook(err)
		}
		return
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
	for _, hook := range d.requestHooks {
		hook(&request)
	}

	if err := d.perform(&request, NewGateway(con)); err != nil {
		return err
	}

	for _, hook := range d.responseHooks {
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
