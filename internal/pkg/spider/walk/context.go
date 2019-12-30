package walk

import (
	"net/url"

	"github.com/muniere/glean/internal/pkg/box"
)

type context struct {
	uri *url.URL
}

func (c *context) dict() box.Dict {
	return box.Dict{
		"uri": c.uri.String(),
	}
}

