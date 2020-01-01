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
	prefix      string
	parallel    int
	concurrency int
	overwrite   bool
	logDir      string
	dryRun      bool
	verbose     bool
}

func assemble(cmd *cobra.Command) {
	cmd.Flags().String("address", rpc.LocalAddr, "Address to bind")
	cmd.Flags().Int("port", rpc.Port, "Port to bind")
	cmd.Flags().String("prefix", "", "Base directory to download files")
	cmd.Flags().Int("parallel", task.Parallel, "The number of workers for download")
	cmd.Flags().Int("concurrency", task.Concurrency, "Concurrency of download tasks per worker")
	cmd.Flags().String("log-dir", "", "Path to log directory")
	cmd.Flags().BoolP("dry-run", "n", false, "Do not perform actions actually")
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

	prefix, err := flags.GetString("prefix")
	if err != nil {
		return nil, err
	}

	parallel, err := flags.GetInt("parallel")
	if err != nil {
		return nil, err
	}

	concurrency, err := flags.GetInt("concurrency")
	if err != nil {
		return nil, err
	}

	logDir, err := flags.GetString("log-dir")
	if err != nil {
		return nil, err
	}

	dryRun, err := flags.GetBool("dry-run")
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
		prefix:      prefix,
		parallel:    parallel,
		concurrency: concurrency,
		logDir:      logDir,
		dryRun:      dryRun,
		verbose:     verbose,
	}

	return opts, nil
}
