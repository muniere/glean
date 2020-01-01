package scrape

import (
	"fmt"
	"io"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func output(w io.Writer, res *rpc.Response) error {
	_, err := fmt.Fprintln(w, jsonic.MustEncode(res))
	return err
}
