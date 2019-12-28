package packet

type Request struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func StatusRequest() Request {
	return Request{
		Action:  "status",
		Payload: StatusPayload{},
	}
}

func LaunchRequest(query string) Request {
	return Request{
		Action:  "launch",
		Payload: LaunchPayload{Query: query},
	}
}

func CancelRequest(query string) Request {
	return Request{
		Action:  "cancel",
		Payload: CancelPayload{Query: query},
	}
}
