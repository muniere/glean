package rpc

import (
	"encoding/json"
	"strings"
)

type Request struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func (r *Request) Encode() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *Request) EncodePretty(indent int) (string, error) {
	spacer := strings.Repeat(" ", indent)
	bytes, err := json.MarshalIndent(r, "", spacer)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (r *Request) DecodePayload(v interface{}) error {
	bytes, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}
