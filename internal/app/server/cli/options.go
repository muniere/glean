package cli

import (
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

type options struct {
	address     string
	port        int
	parallel    int
	concurrency int
	minWidth    int
	maxWidth    int
	minHeight   int
	maxHeight   int
	overwrite   bool
	dataDir     string
	logDir      string
	dryRun      bool
	verbose     bool
}

func assemble(flags *pflag.FlagSet) {
	flags.String("address", rpc.LocalAddr, "Address to bind")
	flags.Int("port", rpc.Port, "Port to bind")
	flags.Int("parallel", task.Parallel, "The number of workers for download")
	flags.Int("concurrency", task.Concurrency, "Concurrency of download tasks per worker")
	flags.Int("min-width", -1, "Minimum width of images")
	flags.Int("max-width", -1, "Maximum width of images")
	flags.Int("min-height", -1, "Minimum height of images")
	flags.Int("max-height", -1, "Maximum height of images")
	flags.String("data-dir", "", "Base directory to download files")
	flags.String("log-dir", "", "Path to log directory")
	flags.BoolP("dry-run", "n", false, "Do not perform actions actually")
	flags.BoolP("verbose", "v", false, "Show verbose messages")
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

	parallel, err := flags.GetInt("parallel")
	if err != nil {
		return nil, err
	}

	concurrency, err := flags.GetInt("concurrency")
	if err != nil {
		return nil, err
	}

	minWidth, err := flags.GetInt("min-width")
	if err != nil {
		return nil, err
	}

	maxWidth, err := flags.GetInt("max-width")
	if err != nil {
		return nil, err
	}

	minHeight, err := flags.GetInt("min-height")
	if err != nil {
		return nil, err
	}

	maxHeight, err := flags.GetInt("max-height")
	if err != nil {
		return nil, err
	}

	dataDir, err := flags.GetString("data-dir")
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
		parallel:    parallel,
		concurrency: concurrency,
		minWidth:    minWidth,
		maxWidth:    maxWidth,
		minHeight:   minHeight,
		maxHeight:   maxHeight,
		dataDir:     dataDir,
		logDir:      logDir,
		dryRun:      dryRun,
		verbose:     verbose,
	}

	return opts, nil
}
