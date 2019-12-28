package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/pkg/agent"
	"github.com/muniere/glean/internal/pkg/defaults"
	"github.com/muniere/glean/internal/pkg/packet"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "status",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	agt := agent.New(defaults.Host, defaults.Port)

	msg := packet.StatusRequest()
	res, err := agt.Submit(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}
