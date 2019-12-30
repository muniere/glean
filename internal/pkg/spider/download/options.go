package download

import (
	"time"
)

type Options struct {
	Prefix      string
	Concurrency int
	Blocking    bool
	Overwrite   bool
	DryRun      bool
	Interval    time.Duration
}

