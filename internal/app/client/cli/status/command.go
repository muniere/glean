package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/muniere/glean/internal/pkg/jsonic"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "status",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	agt := rpc.NewAgent(rpc.RemoteAddr, rpc.Port)

	req := rpc.StatusRequest()
	res, err := agt.Submit(&req)
	if err != nil {
		return err
	}

	var jobs []task.Job
	if err := jsonic.Transcode(res.Payload, &jobs); err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	fmt.Println("id\tURI\tTimestamp")
	for _, job := range jobs {
		fmt.Printf("%d\t%s\t%s\n", job.ID, job.URI, job.Timestamp)
	}

	return nil
}
