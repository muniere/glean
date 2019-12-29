package launch

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "launch",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	var uris []*url.URL
	var errs []string

	for _, arg := range args {
		uri, err := url.Parse(arg)
		if err == nil && len(uri.Scheme) > 0 && len(uri.Host) > 0 {
			uris = append(uris, uri)
		} else {
			errs = append(errs, arg)
		}
	}

	if len(errs) > 0 {
		arg := strings.Join(errs, ", ")
		msg := fmt.Sprintf("values must be valid URLs: %s", arg)
		return errors.New(msg)
	}

	agt := rpc.NewAgent(rpc.RemoteAddr, rpc.Port)

	for _, uri := range uris {
		req := rpc.LaunchRequest(uri)
		res, err := agt.Submit(&req)
		if err != nil {
			return err
		}

		fmt.Println(
			jsonic.MustEncodePretty(res, 4),
		)
	}

	return nil
}
