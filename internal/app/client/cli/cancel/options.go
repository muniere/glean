package cancel

import (
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/app/client/cli/shared"
)

type options struct {
	*shared.Options
}

func assemble(flags *pflag.FlagSet) {
	shared.Assemble(flags)
}

func decode(flags *pflag.FlagSet) (*options, error) {
	base, err := shared.Decode(flags)
	if err != nil {
		return nil, err
	}

	opts := &options{base}

	return opts, nil
}
