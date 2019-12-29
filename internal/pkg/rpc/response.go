package rpc

type Response struct {
	Ok      bool        `json:"ok"`
	Payload interface{} `json:"payload"`
}

