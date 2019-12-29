package rpc

import (
	"net/url"
)

//
// Status
//
type StatusPayload struct{}

func StatusRequest() Request {
	return Request{
		Action:  "status",
		Payload: StatusPayload{},
	}
}

//
// Launch
//
type LaunchPayload struct {
	URI string `json:"uri"`
}

func LaunchRequest(uri *url.URL) Request {
	return Request{
		Action:  "launch",
		Payload: LaunchPayload{URI: uri.String()},
	}
}

//
// Cancel
//
type CancelPayload struct {
	ID int `json:"id"`
}

func CancelRequest(id int) Request {
	return Request{
		Action:  "cancel",
		Payload: CancelPayload{ID: id},
	}
}
