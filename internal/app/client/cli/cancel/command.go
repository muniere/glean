package cancel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "cancel",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	var ids []int
	var errs []string

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err == nil {
			ids = append(ids, id)
		} else {
			errs = append(errs, arg)
		}
	}

	if len(errs) > 0 {
		arg := strings.Join(errs, ", ")
		msg := fmt.Sprintf("values must be ID numbers: %v", arg)
		return errors.New(msg)
	}

	agt := rpc.NewAgent(rpc.RemoteAddr, rpc.Port)

	for _, id := range ids {
		req := rpc.CancelRequest(id)
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
