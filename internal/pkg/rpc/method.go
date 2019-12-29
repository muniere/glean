package rpc

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
	Query string `json:"query"`
}

func LaunchRequest(query string) Request {
	return Request{
		Action:  "launch",
		Payload: LaunchPayload{Query: query},
	}
}

//
// Cancel
//
type CancelPayload struct {
	Query string `json:"query"`
}

func CancelRequest(query string) Request {
	return Request{
		Action:  "cancel",
		Payload: CancelPayload{Query: query},
	}
}
