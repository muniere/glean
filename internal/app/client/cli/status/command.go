package status

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	cli "github.com/muniere/glean/internal/app/client/cli/base"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

func NewCommand() *cobra.Command {
	return assemble(&cobra.Command{
		Use:  "status",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, cmd.Flags())
		},
	})
}

type context struct {
	options optionSet
}

type optionSet struct {
	cli.OptionSet
}

func assemble(cmd *cobra.Command) *cobra.Command {
	return cli.Assemble(cmd)
}

func run(args []string, flags *pflag.FlagSet) error {
	ctx, err := parse(args, flags)
	if err != nil {
		return err
	}

	if err := prepare(ctx); err != nil {
		return err
	}

	agt := rpc.NewAgent(ctx.options.Host, ctx.options.Port)

	req := rpc.NewStatusRequest()
	res, err := agt.Submit(&req)
	if err != nil {
		return err
	}

	return output(os.Stdout, res)
}

func parse(args []string, flags *pflag.FlagSet) (context, error) {
	options, err := decode(flags)
	if err != nil {
		return context{}, err
	}

	ctx := context{
		options: options,
	}

	return ctx, nil
}

func decode(flags *pflag.FlagSet) (optionSet, error) {
	base, err := cli.Decode(flags)
	if err != nil {
		return optionSet{}, err
	}

	opts := optionSet{base}

	return opts, nil
}

func prepare(ctx context) error {
	return cli.Prepare(ctx.options.OptionSet)
}

func output(w io.Writer, res *rpc.Response) error {
	var jobs []task.Job
	if err := jsonic.Transcode(res.Payload, &jobs); err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	if err := puts(w, "ID", "Kind", "URI", "Prefix", "Timestamp"); err != nil {
		return err
	}

	for _, job := range jobs {
		if err := puts(w, string(job.ID), job.Kind, job.URI, job.Prefix, job.Timestamp.String()); err != nil {
			return err
		}
	}

	return nil
}

func puts(w io.Writer, values ...string) error {
	_, err := fmt.Fprintln(w, strings.Join(values, "\t"))
	return err
}
