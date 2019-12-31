package cancel

import (
	"github.com/spf13/pflag"
)

type context struct {
	ids     []int
	options *options
}

func parse(args []string, flags *pflag.FlagSet) (*context, error) {
	ids, err := normalize(args)
	if err != nil {
		return nil, err
	}

	options, err := decode(flags)
	if err != nil {
		return nil, err
	}

	ctx := &context{
		ids:     ids,
		options: options,
	}

	return ctx, nil
}
