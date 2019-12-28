package action

import (
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/pkg/packet"
)

func Uncaught(w *relay.Gateway, req *packet.Request) error {
	return w.Error(nil)
}
