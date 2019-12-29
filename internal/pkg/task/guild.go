package task

import (
	"sync"
	"time"
)

type Guild struct {
	seq     int
	group   *sync.WaitGroup
	workers []*worker
	mutex   *sync.Mutex
}

func NewGuild() *Guild {
	return &Guild{
		seq:     0,
		group:   &sync.WaitGroup{},
		workers: []*worker{},
		mutex:   &sync.Mutex{},
	}
}

func (g *Guild) Spawn(queue *Queue, action Action, recovery Recovery, interval time.Duration) {
	g.mutex.Lock()

	defer g.mutex.Unlock()

	w := &worker{
		id:       g.seq + 1,
		group:    g.group,
		queue:    queue,
		action:   action,
		recovery: recovery,
		interval: interval,
	}

	g.workers = append(g.workers, w)
	g.seq = w.id
}

func (g *Guild) Start() error {
	g.mutex.Lock()

	defer g.mutex.Unlock()

	for _, w := range g.workers {
		if err := w.start(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Guild) Stop() error {
	g.mutex.Lock()

	defer g.mutex.Unlock()

	for _, w := range g.workers {
		if err := w.stop(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Guild) Wait() {
	g.group.Wait()
}
