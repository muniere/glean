package config

import (
	"github.com/spf13/pflag"
)

type context struct {
	*options
}

func parse(args []string, flags *pflag.FlagSet) (*context, error) {
	options, err := decode(flags)
	if err != nil {
		return nil, err
	}

	ctx := &context{
		options: options,
	}

	return ctx, nil
}
