package shared

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Prepare(options *Options) error {
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
