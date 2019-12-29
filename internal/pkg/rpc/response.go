package rpc

import (
	"encoding/json"
	"strings"
)

type Response struct {
	Ok      bool        `json:"ok"`
	Payload interface{} `json:"payload"`
}

func (r *Response) Encode() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *Response) EncodePretty(indent int) (string, error) {
	spacer := strings.Repeat(" ", indent)
	bytes, err := json.MarshalIndent(r, "", spacer)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (r *Response) DecodePayload(v interface{}) error {
	bytes, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}
