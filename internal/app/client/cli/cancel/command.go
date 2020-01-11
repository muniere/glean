package cancel

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/app/client/cli/shared"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func NewCommand() *cobra.Command {
	return assemble(&cobra.Command{
		Use:  "cancel",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, cmd.Flags())
		},
	})
}

type context struct {
	args    argSet
	options optionSet
}

type argSet struct {
	ids []int
}

type optionSet struct {
	shared.OptionSet
}

func assemble(cmd *cobra.Command) *cobra.Command {
	return shared.Assemble(cmd)
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

	for _, id := range ctx.args.ids {
		req := rpc.NewCancelRequest(id)
		res, err := agent.Submit(&req)
		if err != nil {
			return err
		}

		if err := output(os.Stdout, res); err != nil {
			return err
		}
	}

	return nil
}

func parse(args []string, flags *pflag.FlagSet) (context, error) {
	argSet, err := normalize(args)
	if err != nil {
		return context{}, err
	}

	optionSet, err := decode(flags)
	if err != nil {
		return context{}, err
	}

	ctx := context{
		args:    argSet,
		options: optionSet,
	}

	return ctx, nil
}

func normalize(args []string) (argSet, error) {
	var ids []int
	var errs []string

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err == nil {
			ids = append(ids, id)
		} else {
			errs = append(errs, arg)
		}
	}

	if len(errs) > 0 {
		arg := strings.Join(errs, ", ")
		msg := fmt.Sprintf("values must be ID numbers: %v", arg)
		return argSet{}, errors.New(msg)
	}

	return argSet{ids: ids}, nil
}

func decode(flags *pflag.FlagSet) (optionSet, error) {
	base, err := shared.Decode(flags)
	if err != nil {
		return optionSet{}, err
	}

	opts := optionSet{base}

	return opts, nil
}

func prepare(ctx context) error {
	return shared.Prepare(ctx.options.OptionSet)
}

func output(w io.Writer, res *rpc.Response) error {
	_, err := fmt.Fprintln(w, jsonic.MustEncode(res))
	return err
}
