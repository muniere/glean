package task

import (
	"encoding/json"
	"strings"
	"time"
)

type Job struct {
	ID        int       `json:"id"`
	URI       string    `json:"uri"`
	Timestamp time.Time `json:"timestamp"`
}

func (j *Job) Encode() (string, error) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (j *Job) EncodePretty(indent int) (string, error) {
	spacer := strings.Repeat(" ", indent)
	bytes, err := json.MarshalIndent(j, "", spacer)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
