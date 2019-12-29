package action

import (
	"github.com/muniere/glean/internal/app/server/relay"
	"github.com/muniere/glean/internal/app/server/scope"
)

func Uncaught(w *relay.Gateway, req *scope.Context) error {
	return w.Error(nil)
}
