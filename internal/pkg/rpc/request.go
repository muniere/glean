package rpc

type Request struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}
