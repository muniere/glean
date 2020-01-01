package cancel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/client/cli/shared"
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

	assemble(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx, err := parse(args, cmd.Flags())
	if err != nil {
		return err
	}

	if err := shared.Prepare(ctx.options.Options); err != nil {
		return err
	}

	agt := rpc.NewAgent(ctx.options.Host, ctx.options.Port)

	for _, id := range ctx.ids {
		req := rpc.NewCancelRequest(id)
		res, err := agt.Submit(&req)
		if err != nil {
			return err
		}

		fmt.Println(jsonic.MustEncode(res))
	}

	return nil
}
