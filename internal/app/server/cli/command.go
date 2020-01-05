package cli

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/server/pubsub"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/signals"
)

const (
	cmdLogName = "glean.cmd.log"
	outLogName = "glean.out.log"
	errLogName = "glean.err.log"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gleand",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	assemble(cmd)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx, err := parse(args, cmd.Flags())
	if err != nil {
		return err
	}

	err = prepare(ctx.options)
	if err != nil {
		return err
	}

	lumber.Info(box.Dict{
		"module": "root",
		"event":  "start",
		"pid":    os.Getpid(),
	})

	// build
	manager := pubsub.NewManager(translate(ctx.options))

	// start
	err = manager.Start()
	if err != nil {
		lumber.Fatal(box.Dict{
			"module": "root",
			"event":  "start::error",
			"error":  err,
		})
	}

	// wait
	sigs := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
	}
	lumber.Info(box.Dict{
		"module":  "root",
		"event":   "signal::wait",
		"signals": signals.Join(sigs, ", "),
	})

	sig := signals.Wait(sigs...)
	lumber.Info(box.Dict{
		"module": "root",
		"event":  "signal::recv",
		"signal": sig.String(),
	})

	// stop
	err = manager.Stop()
	if err != nil {
		lumber.Fatal(box.Dict{
			"module": "root",
			"event":  "stop::error",
			"error":  err,
		})
	}

	lumber.Info(box.Dict{
		"module": "root",
		"event":  "stop",
		"pid":    os.Getpid(),
	})

	return nil
}

func translate(options *options) pubsub.Config {
	return pubsub.Config{
		Address:     options.address,
		Port:        options.port,
		Prefix:      options.dataDir,
		Parallel:    options.parallel,
		Concurrency: options.concurrency,
		Overwrite:   options.overwrite,
		LogDir:      options.logDir,
		DryRun:      options.dryRun,
		Verbose:     options.verbose,
	}
}
