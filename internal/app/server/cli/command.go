package cli

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/server/pubsub"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/signals"
)

const (
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

	lumber.Debug(box.Dict{
		"module": "root",
		"action": "launch",
		"pid":    os.Getpid(),
	})

	defer lumber.Debug(box.Dict{
		"module": "root",
		"action": "halt",
		"pid":    os.Getpid(),
	})

	// build
	supervisor := pubsub.NewSupervisor(translate(ctx.options))

	// start
	err = supervisor.Start()
	if err != nil {
		log.Fatal(err)
	}

	// wait
	sigs := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
	}
	lumber.Info(box.Dict{
		"module": "root",
		"action": "signal.wait",
		"values": signals.Join(sigs, ", "),
	})

	sig := signals.Wait(sigs...)
	lumber.Info(box.Dict{
		"module": "root",
		"action": "signal.recv",
		"value":  sig.String(),
	})

	// stop
	err = supervisor.Stop()
	if err != nil {
		log.Error(err)
	}

	return nil
}

func translate(options *options) pubsub.Config {
	return pubsub.Config{
		Address:     options.address,
		Port:        options.port,
		Prefix:      options.prefix,
		Parallel:    options.parallel,
		Concurrency: options.concurrency,
		Overwrite:   options.overwrite,
		LogDir:      options.logDir,
		DryRun:      options.dryRun,
		Verbose:     options.verbose,
	}
}
