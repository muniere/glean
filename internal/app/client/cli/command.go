package cli

import (
	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/client/cli/cancel"
	"github.com/muniere/glean/internal/app/client/cli/clutch"
	"github.com/muniere/glean/internal/app/client/cli/config"
	"github.com/muniere/glean/internal/app/client/cli/scrape"
	"github.com/muniere/glean/internal/app/client/cli/status"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "glean",
	}

	cmd.AddCommand(status.NewCommand())
	cmd.AddCommand(scrape.NewCommand())
	cmd.AddCommand(clutch.NewCommand())
	cmd.AddCommand(cancel.NewCommand())
	cmd.AddCommand(config.NewCommand())

	return cmd
}
