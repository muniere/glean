package status

import (
	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/app/client/cli/shared"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "status",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	assemble(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx, err := parse(args, cmd.Flags())
	if err != nil {
		return err
	}

	if err := shared.Prepare(ctx.options.Options); err != nil {
		return err
	}

	agt := rpc.NewAgent(ctx.options.Host, ctx.options.Port)

	req := rpc.NewStatusRequest()
	res, err := agt.Submit(&req)
	if err != nil {
		return err
	}

	var jobs []task.Job
	if err := jsonic.Transcode(res.Payload, &jobs); err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	printLine("ID", "Kind", "URI", "Prefix", "Timestamp")

	for _, job := range jobs {
		printLine(string(job.ID), job.Kind, job.URI, job.Prefix, job.Timestamp.String())
	}

	return nil
}
