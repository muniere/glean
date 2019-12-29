package task

import (
	"time"
)

type Meta struct {
	Worker    int       `json:"worker"`
	Timestamp time.Time `json:"timestamp"`
}
