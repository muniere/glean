package task

import (
	"errors"
	"sync"
	"time"
)

type Action func(Job, Meta) error
type Recovery func(error)

var AlreadyRunning = errors.New("already running")
var NotRunning = errors.New("not running")

type worker struct {
	id       int
	group    *sync.WaitGroup
	queue    *Queue
	action   Action
	recovery Recovery
	interval time.Duration
	halt     func()
}

func (w *worker) start() error {
	if w.halt != nil {
		return AlreadyRunning
	}

	active := true

	w.halt = func() {
		active = false
	}

	w.group.Add(1)

	go func() {
		defer w.group.Done()

		for active {
			job := w.queue.Wait()
			meta := Meta{
				Worker:    w.id,
				Timestamp: time.Now(),
			}

			if err := w.action(job, meta); err != nil {
				w.recovery(err)
				continue
			}

			time.Sleep(w.interval)
		}
	}()

	return nil
}

func (w *worker) stop() error {
	if w.halt == nil {
		return NotRunning
	}

	w.halt()
	w.halt = nil
	return nil
}
