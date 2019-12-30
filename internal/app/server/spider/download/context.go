package download

import (
	"net/url"

	"github.com/muniere/glean/internal/pkg/box"
)

type context struct {
	uri  *url.URL
	path string
}

func (c *context) dict() box.Dict {
	return box.Dict{
		"uri":  c.uri.String(),
		"path": c.path,
	}
}
