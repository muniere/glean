package jsonic

import (
	"encoding/json"
	"strings"
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MustMarshal(v interface{}) []byte {
	bytes, err := Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func MustUnmarshal(data []byte, v interface{}) {
	if err := Unmarshal(data, v); err != nil {
		panic(err)
	}
}

func Encode(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MustEncode(v interface{}) string {
	str, err := Encode(v)
	if err != nil {
		panic(err)
	}
	return str
}

func EncodePretty(v interface{}, indent int) (string, error) {
	spacer := strings.Repeat(" ", indent)
	bytes, err := json.MarshalIndent(v, "", spacer)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MustEncodePretty(v interface{}, indent int) string {
	str, err := EncodePretty(v, indent)
	if err != nil {
		panic(err)
	}
	return str
}

func Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func MustDecode(data []byte, v interface{}) {
	if err := Decode(data, v); err != nil {
		panic(err)
	}
}

func Transcode(v interface{}, u interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, u)
}

func MustTranscode(v interface{}, u interface{}) {
	if err := Transcode(v, u); err != nil {
		panic(err)
	}
}
