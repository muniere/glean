package task

import (
	"errors"
	"sync"
)

type Queue struct {
	seq   int
	jobs  []Job
	mutex *sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		seq:   0,
		jobs:  []Job{},
		mutex: &sync.Mutex{},
	}
}

func (q *Queue) List() []Job {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	return append([]Job{}, q.jobs...)
}

func (q *Queue) Enqueue(query string) error {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	for _, j := range q.jobs {
		if j.Query == query {
			return errors.New("duplicate")
		}
	}

	j := Job{
		ID:    q.seq + 1,
		Query: query,
	}

	q.jobs = append(q.jobs, j)
	q.seq = j.ID

	return nil
}

func (q *Queue) Dequeue() (Job, error) {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	if len(q.jobs) == 0 {
		return Job{}, errors.New("empty")
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job, nil
}

func (q *Queue) Remove(query string) error {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	n := 0
	m := 0
	for _, j := range q.jobs {
		if j.Query == query {
			m++
		} else {
			q.jobs[n] = j
			n++
		}
	}

	if m == 0 {
		return errors.New("not found")
	}

	q.jobs = q.jobs[:n]
	return nil
}
