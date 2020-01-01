package clutch

import (
	"strings"

	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/app/client/cli/shared"
)

type options struct {
	*shared.Options
	Prefix string
}

func assemble(flags *pflag.FlagSet) {
	shared.Assemble(flags)

	flags.StringP("prefix", "p", "", strings.Join([]string{
		"Directory to download files.",
		"Absolute path is resolved as it is.",
		"Relative path is resolved from base directory of glean server.",
	}, "\n"))
}

func decode(flags *pflag.FlagSet) (*options, error) {
	base, err := shared.Decode(flags)
	if err != nil {
		return nil, err
	}

	prefix, err := flags.GetString("prefix")
	if err != nil {
		return nil, err
	}

	opts := &options{base, prefix}

	return opts, nil
}
