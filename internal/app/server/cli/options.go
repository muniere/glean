package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type options struct {
	address     string
	port        int
	concurrency int
	verbose     bool
}

func assemble(cmd *cobra.Command) {
	cmd.Flags().String("address", rpc.LocalAddr, "Address to bind")
	cmd.Flags().Int("port", rpc.Port, "Port to bind")
	cmd.Flags().Int("concurrency", task.Concurrency, "Concurrency of tasks")
	cmd.Flags().BoolP("verbose", "v", false, "Show verbose messages")
}

func decode(flags *pflag.FlagSet) (*options, error) {
	address, err := flags.GetString("address")
	if err != nil {
		return nil, err
	}

	port, err := flags.GetInt("port")
	if err != nil {
		return nil, err
	}

	concurrency, err := flags.GetInt("concurrency")
	if err != nil {
		return nil, err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return nil, err
	}

	opts := &options{
		address:     address,
		port:        port,
		concurrency: concurrency,
		verbose:     verbose,
	}

	return opts, nil
}
