package config

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	. "github.com/muniere/glean/internal/app/client/cli/axiom"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func NewCommand() *cobra.Command {
	return assemble(&cobra.Command{
		Use:  "config",
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
	OptionSet
}

func assemble(cmd *cobra.Command) *cobra.Command {
	return Assemble(cmd)
}

func run(args []string, flags *pflag.FlagSet) error {
	ctx, err := parse(args, flags)
	if err != nil {
		return err
	}

	if err := prepare(ctx); err != nil {
		return err
	}

	agent := rpc.NewAgent(ctx.options.Host, ctx.options.Port)

	req := rpc.NewConfigRequest()
	res, err := agent.Submit(&req)
	if err != nil {
		return err
	}

	return output(os.Stdout, res)
}

func parse(args []string, flags *pflag.FlagSet) (context, error) {
	optionSet, err := decode(flags)
	if err != nil {
		return context{}, err
	}

	ctx := context{
		options: optionSet,
	}

	return ctx, nil
}

func decode(flags *pflag.FlagSet) (optionSet, error) {
	base, err := Decode(flags)
	if err != nil {
		return optionSet{}, err
	}

	opts := optionSet{base}

	return opts, nil
}

func prepare(ctx context) error {
	return Prepare(ctx.options.OptionSet)
}

func output(w io.Writer, res *rpc.Response) error {
	_, err := w.Write(jsonic.MustMarshal(res.Payload))
	return err
}
