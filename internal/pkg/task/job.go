package task

import (
	"time"
)

type Job struct {
	ID        int       `json:"id"`
	URI       string    `json:"uri"`
	Timestamp time.Time `json:"timestamp"`
}
