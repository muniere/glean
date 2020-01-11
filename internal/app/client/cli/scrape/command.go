package scrape

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/app/client/cli/shared"
	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
)

func NewCommand() *cobra.Command {
	return assemble(&cobra.Command{
		Use:  "scrape",
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
	uris []*url.URL
}

type optionSet struct {
	shared.OptionSet
	Prefix string
}

func assemble(cmd *cobra.Command) *cobra.Command {
	flags := cmd.Flags()
	flags.StringP("prefix", "p", "", strings.Join([]string{
		"Directory to download files.",
		"Absolute path is resolved as it is.",
		"Relative path is resolved from base directory of glean server.",
	}, "\n"))

	return shared.Assemble(cmd)
}

func run(args []string, flags *pflag.FlagSet) error {
	ctx, err := parse(args, flags)
	if err != nil {
		return err
	}

	if err := shared.Prepare(ctx.options.OptionSet); err != nil {
		return err
	}

	agt := rpc.NewAgent(ctx.options.Host, ctx.options.Port)

	for _, uri := range ctx.args.uris {
		req := rpc.NewScrapeRequest(uri, ctx.options.Prefix)
		res, err := agt.Submit(&req)
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
	var uris []*url.URL
	var errs []string

	for _, arg := range args {
		uri, err := url.Parse(arg)
		if err == nil && len(uri.Scheme) > 0 && len(uri.Host) > 0 {
			uris = append(uris, uri)
		} else {
			errs = append(errs, arg)
		}
	}

	if len(errs) > 0 {
		arg := strings.Join(errs, ", ")
		msg := fmt.Sprintf("values must be valid URLs: %s", arg)
		return argSet{}, errors.New(msg)
	}

	return argSet{uris: uris}, nil
}

func decode(flags *pflag.FlagSet) (optionSet, error) {
	base, err := shared.Decode(flags)
	if err != nil {
		return optionSet{}, err
	}

	prefix, err := flags.GetString("prefix")
	if err != nil {
		return optionSet{}, err
	}

	opts := optionSet{base, prefix}

	return opts, nil
}

func output(w io.Writer, res *rpc.Response) error {
	_, err := fmt.Fprintln(w, jsonic.MustEncode(res))
	return err
}
