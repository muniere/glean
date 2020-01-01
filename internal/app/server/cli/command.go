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
	"github.com/muniere/glean/internal/pkg/task"
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
	queue := task.NewQueue()

	consumer := pubsub.NewConsumer(queue, pubsub.ConsumerConfig{
		Parallel:    ctx.options.parallel,
		Concurrency: ctx.options.concurrency,
		Prefix:      ctx.options.prefix,
		Overwrite:   ctx.options.overwrite,
		DryRun:      ctx.options.dryRun,
	})

	producer := pubsub.NewProducer(queue, pubsub.ProducerConfig{
		Address: ctx.options.address,
		Port:    ctx.options.port,
	})

	// start
	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = producer.Start()
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
	err = producer.Stop()
	if err != nil {
		log.Error(err)
	}

	err = consumer.Stop()
	if err != nil {
		log.Error(err)
	}

	return nil
}
