package rpc

import (
	"net/url"
)

const (
	Status = "status"
	Scrape = "scrape"
	Cancel = "cancel"
)

//
// Status
//
type StatusPayload struct{}

func StatusRequest() Request {
	return Request{
		Action:  Status,
		Payload: StatusPayload{},
	}
}

//
// Scrape
//
type ScrapePayload struct {
	URI string `json:"uri"`
}

func ScrapeRequest(uri *url.URL) Request {
	return Request{
		Action:  Scrape,
		Payload: ScrapePayload{URI: uri.String()},
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
		Action:  Cancel,
		Payload: CancelPayload{ID: id},
	}
}
