package shared

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/muniere/glean/internal/pkg/rpc"
)

type OptionSet struct {
	Host    string
	Port    int
	Verbose bool
}

func Assemble(cmd *cobra.Command) *cobra.Command {
	flags := cmd.Flags()
	flags.String("host", rpc.RemoteAddr, "Server hostname")
	flags.Int("port", rpc.Port, "Server Port number")
	flags.BoolP("verbose", "v", false, "Show Verbose messages")
	return cmd
}

func Decode(flags *pflag.FlagSet) (OptionSet, error) {
	host, err := flags.GetString("host")
	if err != nil {
		return OptionSet{}, err
	}

	port, err := flags.GetInt("port")
	if err != nil {
		return OptionSet{}, err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return OptionSet{}, err
	}

	opts := OptionSet{
		Host:    host,
		Port:    port,
		Verbose: verbose,
	}

	return opts, nil
}

func Prepare(options OptionSet) error {
	if options.Verbose {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "15:04:05.000",
	})

	return nil
}
