package action

import (
	"github.com/muniere/glean/internal/app/server/scope"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func Uncaught(w *rpc.Gateway, req *scope.Context) error {
	return w.Error(nil)
}
