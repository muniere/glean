package status

import (
	"fmt"
	"io"
	"strings"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

func output(w io.Writer, res *rpc.Response) error {
	var jobs []task.Job
	if err := jsonic.Transcode(res.Payload, &jobs); err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	if err := puts(w, "ID", "Kind", "URI", "Prefix", "Timestamp"); err != nil {
		return err
	}

	for _, job := range jobs {
		if err := puts(w, string(job.ID), job.Kind, job.URI, job.Prefix, job.Timestamp.String()); err != nil {
			return err
		}
	}

	return nil
}

func puts(w io.Writer, values ...string) error {
	_, err := fmt.Fprintln(w, strings.Join(values, "\t"))
	return err
}
