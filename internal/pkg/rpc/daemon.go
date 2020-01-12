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

type Phase struct {
	name string
}

func (p *Phase) String() string {
	return p.name
}

var (
	Accept = Phase{"accept"}
	Handle = Phase{"handle"}
)

type Action interface {
	Invoke(*Request, *Gateway) error
}

type RequestHook interface {
	Invoke(*Request) error
}

type ResponseHook interface {
	Invoke(*Request) error
}

type ErrorHook interface {
	Invoke(error)
}

type Daemon struct {
	Address       string
	Port          int
	Window        int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	listener      net.Listener
	actions       map[string]Action
	fallback      Action
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
		actions:       map[string]Action{},
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
	d.errorHooks[phase] = append(d.subscriptErrorHooks(phase), hook)
}

func (d *Daemon) subscriptErrorHooks(phase Phase) []ErrorHook {
	arr, ok := d.errorHooks[phase]
	if ok {
		return arr
	} else {
		return []ErrorHook{}
	}
}

func (d *Daemon) Register(key string, action Action) {
	d.actions[key] = action
}

func (d *Daemon) Unregister(key string) {
	delete(d.actions, key)
}

func (d *Daemon) RegisterDefault(action Action) {
	d.fallback = action
}

func (d *Daemon) UnregisterDefault() {
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
			err := d.accept()

			if err == nil {
				continue
			} else {
				d.hookError(Accept, err)
				break
			}
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

func (d *Daemon) accept() error {
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
		} else {
			d.hookError(Handle, err)
			return
		}
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
	if err := d.hookRequest(&request); err != nil {
		return err
	}

	if err := d.perform(&request, NewGateway(con)); err != nil {
		return err
	}

	if err := d.hookResponse(&request); err != nil {
		return err
	}

	return nil
}

func (d *Daemon) perform(request *Request, gateway *Gateway) error {
	action, ok := d.actions[request.Action]

	if ok {
		return action.Invoke(request, gateway)
	}

	if d.fallback != nil {
		return d.fallback.Invoke(request, gateway)
	}

	return errors.New(fmt.Sprintf("unsupported aciton: %s", request.Action))
}

func (d *Daemon) hookRequest(req *Request) error {
	for _, hook := range d.requestHooks {
		if err := hook.Invoke(req); err != nil {
			return err
		}
	}
	return nil
}

func (d *Daemon) hookResponse(req *Request) error {
	for _, hook := range d.responseHooks {
		if err := hook.Invoke(req); err != nil {
			return err
		}
	}
	return nil
}

func (d *Daemon) hookError(phase Phase, err error) {
	for _, hook := range d.subscriptErrorHooks(phase) {
		hook.Invoke(err)
	}
}
