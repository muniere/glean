package shared

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/pkg/rpc"
)

type Options struct {
	Host    string
	Port    int
	Verbose bool
}

func Assemble(cmd *cobra.Command) {
	cmd.Flags().String("host", rpc.RemoteAddr, "Server hostname")
	cmd.Flags().Int("port", rpc.Port, "Server Port number")
	cmd.Flags().BoolP("verbose", "v", false, "Show Verbose messages")
}

func Decode(flags *pflag.FlagSet) (*Options, error) {
	host, err := flags.GetString("host")
	if err != nil {
		return nil, err
	}

	port, err := flags.GetInt("port")
	if err != nil {
		return nil, err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return nil, err
	}

	opts := &Options{
		Host:    host,
		Port:    port,
		Verbose: verbose,
	}

	return opts, nil
}
