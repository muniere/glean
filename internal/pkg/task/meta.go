package task

import (
	"encoding/json"
	"strings"
	"time"
)

type Meta struct {
	Worker    int       `json:"worker"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *Meta) Encode() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (m *Meta) EncodePretty(indent int) (string, error) {
	spacer := strings.Repeat(" ", indent)
	bytes, err := json.MarshalIndent(m, "", spacer)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
