package scrape

import (
	"net/url"

	"github.com/spf13/pflag"
)

type context struct {
	uris    []*url.URL
	options *options
}

func parse(args []string, flags *pflag.FlagSet) (*context, error) {
	uris, err := normalize(args)
	if err != nil {
		return nil, err
	}

	options, err := decode(flags)
	if err != nil {
		return nil, err
	}

	ctx := &context{
		uris:    uris,
		options: options,
	}

	return ctx, nil
}
