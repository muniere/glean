package jsonic

import (
	"encoding/json"
	"errors"
)

type (
	array = []json.RawMessage
	dict  = map[string]json.RawMessage
)

var SkipPath = errors.New("skip path")

func Walk(data json.RawMessage, action func(interface{}) error) error {
	if data[0] == '{' {
		var m dict
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		for _, v := range m {
			if err := Walk(v, action); err != nil {
				return err
			}
		}
		return nil
	}

	if data[0] == '[' {
		var a array
		if err := json.Unmarshal(data, &a); err != nil {
			return err
		}
		for _, v := range a {
			if err := Walk(v, action); err != nil {
				return err
			}
		}
		return nil
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if err := action(v); err != nil && err != SkipPath {
		return err
	}
	return nil
}
