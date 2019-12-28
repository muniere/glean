package action

import (
	"encoding/json"

	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/pkg/packet"
)

func Cancel(w *relay.Gateway, req *packet.Request) error {
	bytes, err := json.Marshal(req.Payload)
	if err != nil {
		return err
	}

	var payload packet.CancelPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return err
	}

	return w.Success(nil)
}
