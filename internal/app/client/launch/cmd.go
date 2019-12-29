package launch

import (
	"fmt"

	"github.com/spf13/cobra"

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
	agt := rpc.NewAgent(rpc.Host, rpc.Port)

	for _, query := range args {
		req := rpc.LaunchRequest(query)
		res, err := agt.Submit(&req)
		if err != nil {
			return err
		}

		str, _ := res.EncodePretty(4)
		fmt.Println(str)
	}

	return nil
}
