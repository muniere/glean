package task

import (
	"errors"
	"sync"
	"time"
)

var QueueEmpty = errors.New("queue is empty")
var JobNotFound = errors.New("job not found")
var JobDuplicated = errors.New("job already exists")

type Queue struct {
	seq     int
	jobs    []Job
	mutex   *sync.Mutex
	channel chan bool
}

func NewQueue() *Queue {
	return &Queue{
		seq:     0,
		jobs:    []Job{},
		mutex:   &sync.Mutex{},
		channel: make(chan bool),
	}
}

func (q *Queue) List() []Job {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	return append([]Job{}, q.jobs...)
}

func (q *Queue) Enqueue(uri string) error {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	for _, j := range q.jobs {
		if j.URI == uri {
			return JobDuplicated
		}
	}

	j := Job{
		ID:        q.seq + 1,
		URI:       uri,
		Timestamp: time.Now(),
	}

	q.jobs = append(q.jobs, j)
	q.seq = j.ID

	go func() {
		q.channel <- true
	}()

	return nil
}

func (q *Queue) Dequeue() (Job, error) {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	if len(q.jobs) == 0 {
		return Job{}, QueueEmpty
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	go func() {
		<-q.channel
	}()

	return job, nil
}

func (q *Queue) Wait() Job {
	<-q.channel

	q.mutex.Lock()

	defer q.mutex.Unlock()

	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job
}

func (q *Queue) Remove(id int) error {
	return q.reject(func(job Job) bool {
		return job.ID == id
	})
}

func (q *Queue) Reject(uri string) error {
	return q.reject(func(job Job) bool {
		return job.URI == uri
	})
}

func (q *Queue) reject(test func(Job) bool) error {
	q.mutex.Lock()

	defer q.mutex.Unlock()

	index := 0
	found := false
	for _, job := range q.jobs {
		if found || !test(job) {
			q.jobs[index] = job
			index++
		} else {
			found = true
		}
	}

	if !found {
		return JobNotFound
	}

	q.jobs = q.jobs[:index]

	go func() {
		<-q.channel
	}()

	return nil

}
