package task

import (
	"time"
)

type Job struct {
	ID        int       `json:"id"`
	Kind      string    `json:"kind"`
	URI       string    `json:"uri"`
	Timestamp time.Time `json:"timestamp"`
}
