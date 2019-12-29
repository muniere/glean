package cli

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/server/pubsub"
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
	queue := task.NewQueue()

	consumer := pubsub.NewConsumer(queue, pubsub.ConsumerConfig{
		Concurrency: ctx.options.concurrency,
	})

	producer := pubsub.NewProducer(queue, pubsub.ProducerConfig{
		Address: ctx.options.address,
		Port:    ctx.options.port,
	})

	var err error

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}

	consumer.Wait()
	producer.Wait()

	return nil
}
