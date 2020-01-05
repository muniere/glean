package download

import (
	"time"
)

type Options struct {
	Prefix      string
	Concurrency int
	MinWidth    int
	MaxWidth    int
	MinHeight   int
	MaxHeight   int
	Blocking    bool
	Overwrite   bool
	DryRun      bool
	Interval    time.Duration
}
