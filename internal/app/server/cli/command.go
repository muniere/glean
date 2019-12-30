package cli

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/server/pubsub"
	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/task"
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

	if err := prepare(ctx.options); err != nil {
		return err
	}

	return launch(ctx)
}

func prepare(options *options) error {
	if options.verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "15:04:05.000",
	})

	return nil
}

func launch(ctx *context) error {
	log.Debug(jsonic.MustEncode(box.Dict{
		"module": "root",
		"label":  "launch",
		"pid":    os.Getpid(),
	}))

	defer log.Debug(jsonic.MustEncode(box.Dict{
		"module": "root",
		"label":  "halt",
		"pid":    os.Getpid(),
	}))

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

	// kick
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
	if err := producer.Start(); err != nil {
		log.Fatal(err)
	}

	// wait
	wait(syscall.SIGINT, syscall.SIGTERM)

	if err := producer.Stop(); err != nil {
		log.Error(err)
	}
	if err := consumer.Stop(); err != nil {
		log.Error(err)
	}

	return nil
}

func pinfo() {
}

func wait(sig ...os.Signal) {
	log.Info(jsonic.MustEncode(box.Dict{
		"module": "root",
		"action": "signal",
		"label":  "wait",
		"values": join(sig, ", "),
	}))

	defer log.Info(jsonic.MustEncode(box.Dict{
		"module": "root",
		"action": "signal",
		"label":  "recv",
		"values": join(sig, ", "),
	}))

	ch := make(chan os.Signal)
	signal.Notify(ch, sig...)
	<-ch
}

func join(sig []os.Signal, sep string) string {
	var names []string
	for _, s := range sig {
		names = append(names, s.String())
	}
	return strings.Join(names, sep)
}
