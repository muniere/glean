package rpc

import (
	"net/url"
)

const (
	Status = "status"
	Scrape = "scrape"
	Clutch = "clutch"
	Cancel = "cancel"
	Config = "config"
)

//
// Status
//
type StatusPayload struct{}

func NewStatusRequest() Request {
	return Request{
		Action:  Status,
		Payload: StatusPayload{},
	}
}

//
// Scrape
//
type ScrapePayload struct {
	URI    string `json:"uri"`
	Prefix string `json:"prefix"`
}

func NewScrapeRequest(uri *url.URL, prefix string) Request {
	return Request{
		Action: Scrape,
		Payload: ScrapePayload{
			URI:    uri.String(),
			Prefix: prefix,
		},
	}
}

//
// Clutch
//
type ClutchPayload struct {
	URI    string `json:"uri"`
	Prefix string `json:"prefix"`
}

func NewClutchRequest(uri *url.URL, prefix string) Request {
	return Request{
		Action: Clutch,
		Payload: ClutchPayload{
			URI:    uri.String(),
			Prefix: prefix,
		},
	}
}

//
// Cancel
//
type CancelPayload struct {
	ID int `json:"id"`
}

func NewCancelRequest(id int) Request {
	return Request{
		Action:  Cancel,
		Payload: CancelPayload{ID: id},
	}
}

//
// Config
//
type ConfigPayload struct{}

func NewConfigRequest() Request {
	return Request{
		Action:  Config,
		Payload: ConfigPayload{},
	}
}
