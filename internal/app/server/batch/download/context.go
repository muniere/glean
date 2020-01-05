package download

import (
	"net/url"

	"github.com/muniere/glean/internal/pkg/box"
)

type context struct {
	uri  *url.URL
	temp string
	path string
}

func (c *context) dict() box.Dict {
	dict := box.Dict{}
	if c.uri != nil {
		dict["uri"] = c.uri.String()
	}
	if len(c.temp) > 0 {
		dict["temp"] = c.temp
	}
	if len(c.path) > 0 {
		dict["path"] = c.path
	}
	return dict
}
