package client

import (
	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/client/cancel"
	"github.com/muniere/glean/internal/app/client/launch"
	"github.com/muniere/glean/internal/app/client/status"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "glean",
	}

	cmd.AddCommand(status.NewCommand())
	cmd.AddCommand(launch.NewCommand())
	cmd.AddCommand(cancel.NewCommand())

	return cmd
}
