package config

import (
	"io"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func output(w io.Writer, res *rpc.Response) error {
	_, err := w.Write(jsonic.MustMarshal(res.Payload))
	return err
}
